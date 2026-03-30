package aftersale

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/dao"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

type sAfterSale struct{}

func New() *sAfterSale {
	return &sAfterSale{}
}

func init() {
	service.RegisterAfterSale(New())
}

func (s *sAfterSale) CreateAfterSale(ctx context.Context, req *v1.CreateAfterSaleReq) (*v1.CreateAfterSaleRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetItemNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/item_no are required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	evidenceJSON, _ := json.Marshal(req.GetEvidenceAssetIds())
	afterSaleNo := generateBizNo("AS")
	_, err = dao.AfterSaleCase.Ctx(ctx).Data(do.AfterSaleCase{
		AfterSaleNo:          afterSaleNo,
		OrderNo:              strings.TrimSpace(req.GetOrderNo()),
		SubOrderNo:           strings.TrimSpace(req.GetSubOrderNo()),
		ItemNo:               strings.TrimSpace(req.GetItemNo()),
		UserId:               userID,
		Qty:                  req.GetQty(),
		AfterSaleType:        int(req.GetAfterSaleType()),
		AfterSaleStatus:      int(v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW),
		ApplyRefundAmount:    req.GetApplyRefundAmount(),
		ApprovedRefundAmount: 0,
		ReasonCode:           strings.TrimSpace(req.GetReasonCode()),
		ReasonDesc:           strings.TrimSpace(req.GetReasonDesc()),
		EvidenceAssetIdsJson: string(evidenceJSON),
		BuyerRemark:          strings.TrimSpace(req.GetBuyerRemark()),
		Version:              1,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "create after_sale_case failed")
	}
	row, err := s.getAfterSaleByNo(ctx, afterSaleNo)
	if err != nil {
		return nil, err
	}
	return &v1.CreateAfterSaleRes{AfterSale: toProtoAfterSaleCase(row)}, nil
}

func (s *sAfterSale) CancelAfterSale(ctx context.Context, req *v1.CancelAfterSaleReq) (*v1.CancelAfterSaleRes, error) {
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		return nil, err
	}
	if row.UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "after_sale case does not belong to user")
	}
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	nextVersion := row.Version + 1
	_, err = dao.AfterSaleCase.Ctx(ctx).Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).Data(do.AfterSaleCase{
		AfterSaleStatus:  int(v1.AfterSaleStatus_AFTER_SALE_STATUS_CANCELED),
		CancelReasonCode: strings.TrimSpace(req.GetReasonCode()),
		Version:          nextVersion,
		ClosedAt:         gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "cancel after_sale_case failed")
	}
	return &v1.CancelAfterSaleRes{
		AfterSaleNo:     row.AfterSaleNo,
		AfterSaleStatus: v1.AfterSaleStatus_AFTER_SALE_STATUS_CANCELED,
	}, nil
}

func (s *sAfterSale) GetMyAfterSaleDetail(ctx context.Context, req *v1.GetMyAfterSaleDetailReq) (*v1.GetMyAfterSaleDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		return nil, err
	}
	if row.UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "after_sale case does not belong to user")
	}
	task, _ := s.getRefundTaskByAfterSaleNo(ctx, row.AfterSaleNo)
	return &v1.GetMyAfterSaleDetailRes{
		AfterSale:  toProtoAfterSaleCase(row),
		RefundTask: toProtoRefundTask(task),
	}, nil
}

func (s *sAfterSale) ListMyAfterSales(ctx context.Context, req *v1.ListMyAfterSalesReq) (*v1.ListMyAfterSalesRes, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	model := dao.AfterSaleCase.Ctx(ctx).Where(dao.AfterSaleCase.Columns().UserId, userID).OrderDesc(dao.AfterSaleCase.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		model = model.WhereLT(dao.AfterSaleCase.Columns().Id, cursorID)
	}
	if len(req.GetStatuses()) > 0 {
		statuses := make([]int, 0, len(req.GetStatuses()))
		for _, status := range req.GetStatuses() {
			statuses = append(statuses, int(status))
		}
		model = model.WhereIn(dao.AfterSaleCase.Columns().AfterSaleStatus, statuses)
	}
	var rows []entity.AfterSaleCase
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list my after sales failed")
	}
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	list := make([]*v1.AfterSaleCase, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoAfterSaleCase(&row))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListMyAfterSalesRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

func (s *sAfterSale) ListShopAfterSales(ctx context.Context, req *v1.ListShopAfterSalesReq) (*v1.ListShopAfterSalesRes, error) {
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	model := dao.AfterSaleCase.Ctx(ctx).Where(dao.AfterSaleCase.Columns().ShopNo, req.GetShopNo()).OrderDesc(dao.AfterSaleCase.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		model = model.WhereLT(dao.AfterSaleCase.Columns().Id, cursorID)
	}
	if len(req.GetStatuses()) > 0 {
		statuses := make([]int, 0, len(req.GetStatuses()))
		for _, status := range req.GetStatuses() {
			statuses = append(statuses, int(status))
		}
		model = model.WhereIn(dao.AfterSaleCase.Columns().AfterSaleStatus, statuses)
	}
	var rows []entity.AfterSaleCase
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list shop after sales failed")
	}
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	list := make([]*v1.AfterSaleCase, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoAfterSaleCase(&row))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListShopAfterSalesRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

func (s *sAfterSale) GetShopAfterSaleDetail(ctx context.Context, req *v1.GetShopAfterSaleDetailReq) (*v1.GetShopAfterSaleDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		return nil, err
	}
	task, _ := s.getRefundTaskByAfterSaleNo(ctx, row.AfterSaleNo)
	return &v1.GetShopAfterSaleDetailRes{
		AfterSale:  toProtoAfterSaleCase(row),
		RefundTask: toProtoRefundTask(task),
	}, nil
}

func (s *sAfterSale) ApproveAfterSale(ctx context.Context, req *v1.ApproveAfterSaleReq) (*v1.ApproveAfterSaleRes, error) {
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		return nil, err
	}
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	var task *entity.RefundTask
	err = dao.AfterSaleCase.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		nextVersion := row.Version + 1
		if _, e := tx.Model(dao.AfterSaleCase.Table()).
			Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).
			Data(do.AfterSaleCase{
				AfterSaleStatus:      int(v1.AfterSaleStatus_AFTER_SALE_STATUS_WAIT_REFUND_TASK),
				ApprovedRefundAmount: req.GetApprovedRefundAmount(),
				SellerReply:          strings.TrimSpace(req.GetSellerReply()),
				Version:              nextVersion,
			}).Update(); e != nil {
			return gerror.Wrap(e, "approve after_sale_case failed")
		}
		taskNo := generateBizNo("RF")
		if _, e := tx.Model(dao.RefundTask.Table()).Data(do.RefundTask{
			RefundTaskNo: taskNo,
			AfterSaleNo:  row.AfterSaleNo,
			OrderNo:      row.OrderNo,
			SubOrderNo:   row.SubOrderNo,
			RefundAmount: req.GetApprovedRefundAmount(),
			Status:       int(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
			Version:      1,
		}).Insert(); e != nil {
			return gerror.Wrap(e, "create refund_task failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	row, _ = s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	task, _ = s.getRefundTaskByAfterSaleNo(ctx, req.GetAfterSaleNo())
	return &v1.ApproveAfterSaleRes{
		AfterSale:  toProtoAfterSaleCase(row),
		RefundTask: toProtoRefundTask(task),
	}, nil
}

func (s *sAfterSale) RejectAfterSale(ctx context.Context, req *v1.RejectAfterSaleReq) (*v1.RejectAfterSaleRes, error) {
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		return nil, err
	}
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	nextVersion := row.Version + 1
	_, err = dao.AfterSaleCase.Ctx(ctx).
		Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).
		Data(do.AfterSaleCase{
			AfterSaleStatus:  int(v1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED),
			RejectReasonCode: int(req.GetRejectReasonCode()),
			SellerReply:      strings.TrimSpace(req.GetSellerReply()),
			Version:          nextVersion,
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "reject after_sale_case failed")
	}
	return &v1.RejectAfterSaleRes{
		AfterSaleNo:     row.AfterSaleNo,
		AfterSaleStatus: v1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED,
	}, nil
}

func (s *sAfterSale) ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (*v1.ExecuteRefundTaskRes, error) {
	if req == nil || strings.TrimSpace(req.GetRefundTaskNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_task_no is required")
	}
	task, err := s.getRefundTaskByNo(ctx, req.GetRefundTaskNo())
	if err != nil {
		return nil, err
	}
	if task.Status == uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED) {
		return &v1.ExecuteRefundTaskRes{
			Task:            toProtoRefundTask(task),
			RefundSucceeded: true,
			ShouldRetry:     false,
		}, nil
	}
	err = dao.RefundTask.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := tx.Model(dao.RefundTask.Table()).
			Where(dao.RefundTask.Columns().RefundTaskNo, task.RefundTaskNo).
			Data(do.RefundTask{
				Status:      int(v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED),
				RetryCount:  task.RetryCount + 1,
				NextRetryAt: nil,
			}).Update(); e != nil {
			return gerror.Wrap(e, "mark refund_task succeeded failed")
		}
		if _, e := tx.Model(dao.AfterSaleCase.Table()).
			Where(dao.AfterSaleCase.Columns().AfterSaleNo, task.AfterSaleNo).
			Data(do.AfterSaleCase{
				AfterSaleStatus: int(v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED),
				ClosedAt:        gtime.Now(),
			}).Update(); e != nil {
			return gerror.Wrap(e, "mark after_sale_case refunded failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	task, _ = s.getRefundTaskByNo(ctx, req.GetRefundTaskNo())
	return &v1.ExecuteRefundTaskRes{
		Task:            toProtoRefundTask(task),
		RefundSucceeded: true,
		ShouldRetry:     false,
	}, nil
}

func (s *sAfterSale) RetryRefundTask(ctx context.Context, req *v1.RetryRefundTaskReq) (*v1.RetryRefundTaskRes, error) {
	if req == nil || strings.TrimSpace(req.GetRefundTaskNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_task_no is required")
	}
	_, err := dao.RefundTask.Ctx(ctx).
		Where(dao.RefundTask.Columns().RefundTaskNo, req.GetRefundTaskNo()).
		Data(do.RefundTask{
			Status:           int(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
			NextRetryAt:      gtime.Now(),
			LastErrorCode:    strings.TrimSpace(req.GetReasonCode()),
			LastErrorMessage: "retry requested",
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "retry refund task failed")
	}
	task, err := s.getRefundTaskByNo(ctx, req.GetRefundTaskNo())
	if err != nil {
		return nil, err
	}
	return &v1.RetryRefundTaskRes{Task: toProtoRefundTask(task)}, nil
}

func (s *sAfterSale) GetAfterSaleSnapshotByNo(ctx context.Context, req *v1.GetAfterSaleSnapshotByNoReq) (*v1.GetAfterSaleSnapshotByNoRes, error) {
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no is required")
	}
	row, err := s.getAfterSaleByNo(ctx, req.GetAfterSaleNo())
	if err != nil {
		return nil, err
	}
	task, _ := s.getRefundTaskByAfterSaleNo(ctx, req.GetAfterSaleNo())
	return &v1.GetAfterSaleSnapshotByNoRes{
		AfterSale:  toProtoAfterSaleCase(row),
		RefundTask: toProtoRefundTask(task),
	}, nil
}

func (s *sAfterSale) getAfterSaleByNo(ctx context.Context, afterSaleNo string) (*entity.AfterSaleCase, error) {
	var row entity.AfterSaleCase
	if err := dao.AfterSaleCase.Ctx(ctx).Where(dao.AfterSaleCase.Columns().AfterSaleNo, afterSaleNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query after_sale_case failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "after_sale case not found")
	}
	return &row, nil
}

func (s *sAfterSale) getRefundTaskByNo(ctx context.Context, refundTaskNo string) (*entity.RefundTask, error) {
	var row entity.RefundTask
	if err := dao.RefundTask.Ctx(ctx).Where(dao.RefundTask.Columns().RefundTaskNo, refundTaskNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query refund_task failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "refund task not found")
	}
	return &row, nil
}

func (s *sAfterSale) getRefundTaskByAfterSaleNo(ctx context.Context, afterSaleNo string) (*entity.RefundTask, error) {
	var row entity.RefundTask
	if err := dao.RefundTask.Ctx(ctx).Where(dao.RefundTask.Columns().AfterSaleNo, afterSaleNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query refund task by after_sale_no failed")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

func toProtoAfterSaleCase(row *entity.AfterSaleCase) *v1.AfterSaleCase {
	if row == nil {
		return nil
	}
	var evidence []uint64
	_ = json.Unmarshal([]byte(row.EvidenceAssetIdsJson), &evidence)
	return &v1.AfterSaleCase{
		AfterSaleNo:          row.AfterSaleNo,
		OrderNo:              row.OrderNo,
		SubOrderNo:           row.SubOrderNo,
		ItemNo:               row.ItemNo,
		UserId:               row.UserId,
		ShopNo:               row.ShopNo,
		SpuNo:                row.SpuNo,
		SkuNo:                row.SkuNo,
		Qty:                  uint32(row.Qty),
		AfterSaleType:        v1.AfterSaleType(row.AfterSaleType),
		AfterSaleStatus:      v1.AfterSaleStatus(row.AfterSaleStatus),
		ApplyRefundAmount:    row.ApplyRefundAmount,
		ApprovedRefundAmount: row.ApprovedRefundAmount,
		ReasonCode:           row.ReasonCode,
		ReasonDesc:           row.ReasonDesc,
		EvidenceAssetIds:     evidence,
		BuyerRemark:          row.BuyerRemark,
		SellerReply:          row.SellerReply,
		RejectReasonCode:     v1.RejectReasonCode(row.RejectReasonCode),
		Version:              row.Version,
		CreatedAt:            toProtoTs(row.CreatedAt),
		UpdatedAt:            toProtoTs(row.UpdatedAt),
		ClosedAt:             toProtoTs(row.ClosedAt),
		CancelReasonCode:     row.CancelReasonCode,
	}
}

func toProtoRefundTask(row *entity.RefundTask) *v1.RefundTask {
	if row == nil {
		return nil
	}
	return &v1.RefundTask{
		RefundTaskNo:     row.RefundTaskNo,
		AfterSaleNo:      row.AfterSaleNo,
		OrderNo:          row.OrderNo,
		SubOrderNo:       row.SubOrderNo,
		PayNo:            row.PayNo,
		RefundAmount:     row.RefundAmount,
		Status:           v1.RefundTaskStatus(row.Status),
		RetryCount:       uint32(row.RetryCount),
		NextRetryAt:      toProtoTs(row.NextRetryAt),
		LastErrorCode:    row.LastErrorCode,
		LastErrorMessage: row.LastErrorMessage,
		CreatedAt:        toProtoTs(row.CreatedAt),
		UpdatedAt:        toProtoTs(row.UpdatedAt),
	}
}

func userIDFromContext(ctx context.Context) (uint64, error) {
	if r := g.RequestFromCtx(ctx); r != nil {
		for _, key := range []string{"x-user-id", "X-User-Id", "user_id", "uid"} {
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				uid, err := strconv.ParseUint(raw, 10, 64)
				if err == nil && uid > 0 {
					return uid, nil
				}
			}
		}
	}
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"x-user-id", "user_id", "uid", "userid"} {
		if val := md.Get(key); val != nil {
			uid := gconv.Uint64(val)
			if uid > 0 {
				return uid, nil
			}
		}
	}
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "missing x-user-id")
}

func normalizePageSize(reqSize int32) int {
	size := int(reqSize)
	if size <= 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return size
}

func parseCursor(cursor string) (uint64, error) {
	cursor = strings.TrimSpace(cursor)
	if cursor == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(cursor, 10, 64)
	if err != nil {
		return 0, gerror.WrapCode(gcode.CodeInvalidParameter, err, "next_cursor must be uint64")
	}
	return id, nil
}

func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
