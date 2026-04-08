package aftersale

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/dao"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func (s *sAfterSale) replayRefundBatchIdempotency(ctx context.Context, userID uint64, key, action string) (bool, string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return false, "", nil
	}
	model := dao.AftersaleIdempotency.Ctx(ctx).
		Where(dao.AftersaleIdempotency.Columns().IdempotencyKey, key).
		Where(dao.AftersaleIdempotency.Columns().ActionCode, action)
	if userID > 0 {
		model = model.Where(dao.AftersaleIdempotency.Columns().UserId, userID)
	}
	var row entity.AftersaleIdempotency
	if err := model.Scan(&row); err != nil {
		if isRefundBatchIdempotencyMissError(err) {
			return false, "", nil
		}
		return false, "", gerror.Wrap(err, "query refund batch idempotency failed")
	}
	if row.Id == 0 {
		return false, "", nil
	}
	if strings.TrimSpace(row.ResourceNo) != "" {
		return true, row.ResourceNo, nil
	}
	var payload struct {
		RefundBatchNo string `json:"refund_batch_no"`
	}
	_ = json.Unmarshal([]byte(row.ResponseJson), &payload)
	return true, payload.RefundBatchNo, nil
}

func isRefundBatchIdempotencyMissError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}
	if strings.Contains(strings.ToLower(err.Error()), sql.ErrNoRows.Error()) {
		return true
	}
	return false
}

func (s *sAfterSale) saveRefundBatchIdempotencyTx(ctx context.Context, tx gdb.TX, userID uint64, key, action, refundBatchNo string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}
	payload, _ := json.Marshal(map[string]any{
		"refund_batch_no": refundBatchNo,
	})
	_, err := tx.Model(dao.AftersaleIdempotency.Table()).Data(do.AftersaleIdempotency{
		UserId:         userID,
		IdempotencyKey: key,
		ActionCode:     action,
		ResourceNo:     refundBatchNo,
		Status:         1,
		ResponseJson:   string(payload),
		ExpireAt:       gtime.New(time.Now().Add(24 * time.Hour)),
	}).InsertIgnore()
	if err != nil {
		return gerror.Wrap(err, "save refund batch idempotency failed")
	}
	return nil
}

func (s *sAfterSale) upsertLocalRefundTaskTx(ctx context.Context, tx gdb.TX, row *entity.AfterSaleCase, refundTaskNo string) error {
	insertRes, err := tx.Model(dao.RefundTask.Table()).Data(do.RefundTask{
		RefundTaskNo: refundTaskNo,
		AfterSaleNo:  row.AfterSaleNo,
		OrderNo:      row.OrderNo,
		SubOrderNo:   row.SubOrderNo,
		PayNo:        row.PaymentNo,
		RefundAmount: row.ApplyRefundAmount,
		Status:       uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
		Version:      1,
	}).InsertIgnore()
	if err != nil {
		return gerror.Wrap(err, "create local refund task failed")
	}
	if insertRes != nil {
		if rows, _ := insertRes.RowsAffected(); rows > 0 {
			return nil
		}
	}
	_, err = tx.Model(dao.RefundTask.Table()).
		Where(dao.RefundTask.Columns().AfterSaleNo, row.AfterSaleNo).
		Data(do.RefundTask{
			RefundTaskNo:           refundTaskNo,
			OrderNo:                row.OrderNo,
			SubOrderNo:             row.SubOrderNo,
			PayNo:                  row.PaymentNo,
			RefundAmount:           row.ApplyRefundAmount,
			Status:                 uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
			LastErrorCode:          "",
			LastErrorMessage:       "",
			PointsReturnAmount:     0,
			PointsReverseAmount:    0,
			PointsCashOffsetAmount: 0,
			FinalCashRefundAmount:  0,
			AccountDebtAfter:       0,
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "update local refund task failed")
	}
	return nil
}

func (s *sAfterSale) markLocalRefundTaskFailed(ctx context.Context, task *entity.RefundTask, errorCode, errorMessage string) error {
	return dao.RefundTask.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.RefundTask.Table()).
			Where(dao.RefundTask.Columns().RefundTaskNo, task.RefundTaskNo).
			Data(do.RefundTask{
				Status:           uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_FAILED),
				NextRetryAt:      gtime.New(time.Now().Add(1 * time.Minute)),
				LastErrorCode:    strings.TrimSpace(errorCode),
				LastErrorMessage: strings.TrimSpace(errorMessage),
			}).Update(); err != nil {
			return gerror.Wrap(err, "mark local refund task failed")
		}
		if _, err := tx.Model(dao.AfterSaleCase.Table()).
			Where(dao.AfterSaleCase.Columns().AfterSaleNo, task.AfterSaleNo).
			Data(do.AfterSaleCase{
				AfterSaleStatus: uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUND_PROCESSING),
			}).Update(); err != nil {
			return gerror.Wrap(err, "mark after sale refund processing failed")
		}
		return nil
	})
}

// ProcessExpiredRefundReviewBatches 扫描超时未审核的退款批次，自动通过并沿用人工审批链路。
func ProcessExpiredRefundReviewBatches(ctx context.Context, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 20
	}
	var rows []entity.AfterSaleCase
	if err := dao.AfterSaleCase.Ctx(ctx).
		Where(dao.AfterSaleCase.Columns().AfterSaleStatus, uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW)).
		WhereLTE(dao.AfterSaleCase.Columns().ReviewDeadlineAt, gtime.Now()).
		OrderAsc(dao.AfterSaleCase.Columns().ReviewDeadlineAt).
		Limit(batchSize).
		Scan(&rows); err != nil {
		return gerror.Wrap(err, "query expired refund review batches failed")
	}
	if len(rows) == 0 {
		return nil
	}
	svc := New()
	seen := make(map[string]string)
	for _, row := range rows {
		if strings.TrimSpace(row.RefundBatchNo) == "" {
			continue
		}
		if _, ok := seen[row.RefundBatchNo]; !ok {
			seen[row.RefundBatchNo] = row.ShopNo
		}
	}
	var lastErr error
	for refundBatchNo, shopNo := range seen {
		if _, err := svc.ApproveRefundBatch(ctx, &v1.ApproveRefundBatchReq{
			RefundBatchNo:  refundBatchNo,
			ShopNo:         shopNo,
			SellerReply:    "超时自动通过退款",
			IdempotencyKey: fmt.Sprintf("auto_approve_%s", refundBatchNo),
		}); err != nil {
			lastErr = err
			g.Log().Errorf(ctx, "[aftersale-svc] auto approve refund batch failed, refund_batch_no=%s, err=%+v", refundBatchNo, err)
		}
	}
	return lastErr
}
