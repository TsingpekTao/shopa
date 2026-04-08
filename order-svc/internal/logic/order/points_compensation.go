package order

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/TsingpekTao/shopa/order-svc/internal/dao"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	pointsCompTaskStatusPending = "PENDING"
	pointsCompTaskStatusSuccess = "SUCCESS"
	pointsCompTaskStatusFailed  = "FAILED"
	pointsCompTaskStatusDead    = "DEAD"
)

type pointsCompensationWorkerConf struct {
	Enabled         bool
	IntervalSeconds int
	BatchSize       int
	MaxRetry        int
}

// enqueuePointsCancelCompensationTask 写入待执行的积分补偿任务。
func (s *sOrder) enqueuePointsCancelCompensationTask(ctx context.Context, userID uint64, orderNo, reservationNo, lastError string) error {
	if strings.TrimSpace(reservationNo) == "" {
		return nil
	}
	now := gtime.Now()
	// 将失败原因裁剪后写入 LastError，供后续 worker 判断重试背景。
	_, err := dao.OrderPointsCompensationTask.Ctx(ctx).Data(do.OrderPointsCompensationTask{
		TaskNo:              generateBizNo("OPC"),
		OrderNo:             orderNo,
		UserId:              userID,
		PointsReservationNo: reservationNo,
		ActionCode:          pointsCompensationActionCancel,
		TaskStatus:          pointsCompTaskStatusPending,
		RetryCount:          0,
		NextRetryAt:         now,
		LastError:           trimForColumn(lastError, 1000),
	}).Insert()
	return gerror.Wrap(err, "insert order_points_compensation_task failed")
}

// enqueuePointsGrantCompensationTask 写入待重试的赠积分补偿任务。
func (s *sOrder) enqueuePointsGrantCompensationTask(ctx context.Context, userID uint64, orderNo, lastError string) error {
	now := gtime.Now()
	_, err := dao.OrderPointsCompensationTask.Ctx(ctx).Data(do.OrderPointsCompensationTask{
		TaskNo:              generateBizNo("OPC"),
		OrderNo:             orderNo,
		UserId:              userID,
		PointsReservationNo: "",
		ActionCode:          pointsCompensationActionGrant,
		TaskStatus:          pointsCompTaskStatusPending,
		RetryCount:          0,
		NextRetryAt:         now,
		LastError:           trimForColumn(lastError, 1000),
	}).Insert()
	return gerror.Wrap(err, "insert points grant compensation task failed")
}

// compensateGrantPointsAfterOrderCompleted 订单完成后，若赠积分失败则降级为补偿任务。
func (s *sOrder) compensateGrantPointsAfterOrderCompleted(ctx context.Context, userID uint64, orderNo string, causeErr error) {
	if userID == 0 || strings.TrimSpace(orderNo) == "" {
		return
	}
	if enqueueErr := s.enqueuePointsGrantCompensationTask(ctx, userID, orderNo, causeErr.Error()); enqueueErr != nil {
		g.Log().Errorf(ctx, "enqueue grant points compensation failed, order_no=%s err=%+v", orderNo, enqueueErr)
		return
	}
	g.Log().Warningf(ctx, "grant points failed, fallback compensation task inserted, order_no=%s err=%+v", orderNo, causeErr)
}

// compensateLockedPointsAfterCreateFailure 在下单失败后补偿释放已锁积分。
func (s *sOrder) compensateLockedPointsAfterCreateFailure(ctx context.Context, userID uint64, orderNo, reservationNo string, causeErr error) error {
	if strings.TrimSpace(reservationNo) == "" {
		return nil
	}
	// 先尝试同步取消；若失败则降级为本地补偿任务，避免下游不可用时积分长期占用。
	reasonCode := "ORDER_CREATE_FAILED"
	cancelErr := s.cancelLockedPoints(ctx, userID, orderNo, reservationNo, reasonCode, buildPointsCancelIdempotencyKey(orderNo, reasonCode))
	if cancelErr == nil {
		return nil
	}
	if enqueueErr := s.enqueuePointsCancelCompensationTask(ctx, userID, orderNo, reservationNo, cancelErr.Error()); enqueueErr != nil {
		g.Log().Errorf(ctx, "enqueue points compensation failed, order_no=%s reservation_no=%s err=%+v", orderNo, reservationNo, enqueueErr)
		return gerror.Wrapf(cancelErr, "cancel locked points failed and enqueue task failed, create_err=%v", causeErr)
	}
	g.Log().Warningf(ctx, "sync cancel points failed, fallback task inserted, order_no=%s reservation_no=%s err=%+v", orderNo, reservationNo, cancelErr)
	return nil
}

// compensateLockedPointsAfterOrderClose 在订单关闭后补偿释放积分锁定。
func (s *sOrder) compensateLockedPointsAfterOrderClose(ctx context.Context, userID uint64, orderNo, reservationNo, reasonCode string) {
	if strings.TrimSpace(reservationNo) == "" {
		return
	}
	// 先做同步取消；失败时写入补偿任务，依赖 worker 重试，保证锁定积分最终返还。
	if err := s.cancelLockedPoints(ctx, userID, orderNo, reservationNo, reasonCode, buildPointsCancelIdempotencyKey(orderNo, reasonCode)); err != nil {
		g.Log().Warningf(ctx, "sync cancel locked points failed, order_no=%s reservation_no=%s err=%+v", orderNo, reservationNo, err)
		if enqueueErr := s.enqueuePointsCancelCompensationTask(ctx, userID, orderNo, reservationNo, err.Error()); enqueueErr != nil {
			g.Log().Errorf(ctx, "enqueue cancel points compensation failed, order_no=%s reservation_no=%s err=%+v", orderNo, reservationNo, enqueueErr)
		}
	}
}

// ProcessPendingPointsCompensationTasks 对外暴露补偿任务处理入口。
func ProcessPendingPointsCompensationTasks(ctx context.Context, limit int) error {
	return New().processPendingPointsCompensationTasks(ctx, limit)
}

// processPendingPointsCompensationTasks 扫描并处理待执行的补偿任务。
func (s *sOrder) processPendingPointsCompensationTasks(ctx context.Context, limit int) error {
	if limit <= 0 {
		limit = loadPointsCompensationWorkerConf(ctx).BatchSize
	}
	var rows []*entity.OrderPointsCompensationTask
	cols := dao.OrderPointsCompensationTask.Columns()
	if err := dao.OrderPointsCompensationTask.Ctx(ctx).
		WhereIn(cols.TaskStatus, []string{pointsCompTaskStatusPending, pointsCompTaskStatusFailed}).
		WhereLTE(cols.NextRetryAt, gtime.Now()).
		OrderAsc(cols.Id).
		Limit(limit).
		Scan(&rows); err != nil {
		return gerror.Wrap(err, "query pending points compensation task failed")
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		// 逐条处理；失败时在 handlePointsCompensationTask 内部记录日志并更新重试状态。
		if err := s.handlePointsCompensationTask(ctx, row); err != nil {
			g.Log().Errorf(ctx, "handle points compensation task failed, task_no=%s err=%+v", row.TaskNo, err)
		}
	}
	return nil
}

// handlePointsCompensationTask 执行单条积分补偿任务。
func (s *sOrder) handlePointsCompensationTask(ctx context.Context, row *entity.OrderPointsCompensationTask) error {
	switch strings.TrimSpace(row.ActionCode) {
	case pointsCompensationActionCancel:
		// 当前仅支持取消 reservation；后续若新增赠分补偿，可在此分支扩展。
		err := s.cancelLockedPoints(ctx, row.UserId, row.OrderNo, row.PointsReservationNo, "SYSTEM_COMPENSATION", buildPointsCancelIdempotencyKey(row.OrderNo, "SYSTEM_COMPENSATION"))
		return s.finishPointsCompensationTask(ctx, row, err)
	case pointsCompensationActionGrant:
		mainRow, err := s.getOrderMainByNo(ctx, row.OrderNo)
		if err != nil {
			return s.finishPointsCompensationTask(ctx, row, err)
		}
		if mainRow == nil || strings.TrimSpace(mainRow.OrderNo) == "" {
			return s.finishPointsCompensationTask(ctx, row, gerror.New("order not found"))
		}
		if err := s.grantPointsByOrderCompleted(ctx, row.UserId, row.OrderNo, mainRow.PaidAmount, buildPointsGrantIdempotencyKey(row.OrderNo)); err != nil {
			return s.finishPointsCompensationTask(ctx, row, err)
		}
		return s.finishPointsCompensationTask(ctx, row, nil)
	default:
		return s.finishPointsCompensationTask(ctx, row, gerror.Newf("unsupported compensation action: %s", row.ActionCode))
	}
}

// finishPointsCompensationTask 根据执行结果更新任务状态。
func (s *sOrder) finishPointsCompensationTask(ctx context.Context, row *entity.OrderPointsCompensationTask, execErr error) error {
	if row == nil {
		return nil
	}
	cols := dao.OrderPointsCompensationTask.Columns()
	if execErr == nil {
		// 成功时标记为 SUCCESS，并清空错误信息，避免重复重试。
		_, err := dao.OrderPointsCompensationTask.Ctx(ctx).
			Where(cols.Id, row.Id).
			Data(do.OrderPointsCompensationTask{
				TaskStatus: pointsCompTaskStatusSuccess,
				LastError:  "",
			}).Update()
		return gerror.Wrap(err, "mark points compensation success failed")
	}
	conf := loadPointsCompensationWorkerConf(ctx)
	retryCount := int(row.RetryCount) + 1
	nextStatus := pointsCompTaskStatusFailed
	if retryCount >= conf.MaxRetry {
		nextStatus = pointsCompTaskStatusDead
	}
	// 失败时更新重试次数和下次执行时间，采用平方退避避免频繁打爆下游。
	_, err := dao.OrderPointsCompensationTask.Ctx(ctx).
		Where(cols.Id, row.Id).
		Data(do.OrderPointsCompensationTask{
			TaskStatus:  nextStatus,
			RetryCount:  retryCount,
			NextRetryAt: nextPointsRetryTime(retryCount),
			LastError:   trimForColumn(execErr.Error(), 1000),
		}).Update()
	return gerror.Wrap(err, "update points compensation retry state failed")
}

// loadPointsCompensationWorkerConf 读取积分补偿 worker 配置。
func loadPointsCompensationWorkerConf(ctx context.Context) pointsCompensationWorkerConf {
	var (
		interval = g.Cfg().MustGet(ctx, "worker.pointsCompensation.intervalSeconds", 10).Int()
		batch    = g.Cfg().MustGet(ctx, "worker.pointsCompensation.batchSize", 20).Int()
		maxRetry = g.Cfg().MustGet(ctx, "worker.pointsCompensation.maxRetry", 12).Int()
	)
	if interval <= 0 {
		interval = 10
	}
	if batch <= 0 {
		batch = 20
	}
	if maxRetry <= 0 {
		maxRetry = 12
	}
	// Enabled 从配置读取，允许线上通过开关暂停补偿 worker。
	return pointsCompensationWorkerConf{
		Enabled:         g.Cfg().MustGet(ctx, "worker.pointsCompensation.enabled", true).Bool(),
		IntervalSeconds: interval,
		BatchSize:       batch,
		MaxRetry:        maxRetry,
	}
}

// nextPointsRetryTime 计算下一次补偿重试时间。
func nextPointsRetryTime(retryCount int) *gtime.Time {
	backoff := time.Duration(retryCount*retryCount) * time.Minute
	if backoff > 30*time.Minute {
		backoff = 30 * time.Minute
	}
	return gtime.NewFromTime(time.Now().Add(backoff))
}

// trimForColumn 将字符串裁剪到数据库字段允许长度。
func trimForColumn(in string, maxLen int) string {
	in = strings.TrimSpace(in)
	if maxLen <= 0 || len(in) <= maxLen {
		return in
	}
	return in[:maxLen]
}

// appendOperateLogTx 在事务内追加订单操作日志。
func appendOperateLogTx(ctx context.Context, tx gdb.TX, orderNo, subOrderNo, actionCode, beforeStatus, afterStatus string, detail any) {
	if tx == nil || strings.TrimSpace(orderNo) == "" {
		return
	}
	// detail 可能较大；序列化失败时不阻断主流程，仅忽略详情。
	detailJSON := ""
	if detail != nil {
		if encoded, err := gjsonEncode(detail); err == nil {
			detailJSON = encoded
		}
	}
	_, _ = tx.Model(dao.OrderOperateLog.Table()).Data(do.OrderOperateLog{
		OrderNo:          orderNo,
		SubOrderNo:       subOrderNo,
		OperatorUserId:   uint64(0),
		OperatorTypeCode: "SYSTEM",
		ActionCode:       actionCode,
		BeforeStatus:     beforeStatus,
		AfterStatus:      afterStatus,
		DetailJson:       detailJSON,
	}).Insert()
	_ = ctx
}

// gjsonEncode 将任意结构编码为 JSON 字符串。
func gjsonEncode(v any) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
