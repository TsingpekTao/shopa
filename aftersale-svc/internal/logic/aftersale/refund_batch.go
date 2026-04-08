package aftersale

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	orderv1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	paymentv1 "github.com/TsingpekTao/shopa/payment-svc/api/v1"

	v1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/dao"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// refundScopeSubOrder 表示 V1 虽然允许从商品入口发起，但后台统一按子单级退款处理。
	refundScopeSubOrder = "SUB_ORDER"
	// refundReviewWindow 定义卖家审核窗口，超时后 worker 会自动同意并继续退款。
	refundReviewWindow = 24 * time.Hour

	refundActionApply   = "refund_batch_apply"
	refundActionCancel  = "refund_batch_cancel"
	refundActionApprove = "refund_batch_approve"
	refundActionReject  = "refund_batch_reject"
)

type refundBatchCursorRow struct {
	CursorID      uint64 `json:"cursor_id"`
	RefundBatchNo string `json:"refund_batch_no"`
}

type refundBatchApplyCase struct {
	afterSaleNo       string
	orderNo           string
	subOrderNo        string
	itemNo            string
	selectedItemNos   []string
	snapshot          *orderv1.RefundSubOrderSnapshot
	applyRefundAmount uint64
}

// ApplyRefundBatch 创建买家预发货退款批次，并先把目标子单冻结到退款中，避免卖家继续发货。
func (s *sAfterSale) ApplyRefundBatch(ctx context.Context, req *v1.ApplyRefundBatchReq) (*v1.ApplyRefundBatchRes, error) {
	if req == nil || len(req.GetTargets()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "targets are required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if hit, batchNo, err := s.replayRefundBatchIdempotency(ctx, userID, strings.TrimSpace(req.GetIdempotencyKey()), refundActionApply); err != nil {
		return nil, err
	} else if hit {
		detail, innerErr := s.buildRefundBatchDetail(ctx, batchNo)
		if innerErr != nil {
			return nil, innerErr
		}
		return &v1.ApplyRefundBatchRes{Detail: detail}, nil
	}

	// 先做目标去重，避免前端重复提交同一子单导致创建多张售后单。
	uniqueTargets := make(map[string]*v1.ApplyRefundTarget)
	for _, target := range req.GetTargets() {
		if target == nil || strings.TrimSpace(target.GetOrderNo()) == "" || strings.TrimSpace(target.GetSubOrderNo()) == "" {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no are required")
		}
		key := strings.TrimSpace(target.GetSubOrderNo())
		uniqueTargets[key] = target
	}

	applyCases := make([]refundBatchApplyCase, 0, len(uniqueTargets))
	var (
		batchOrderNo string
		batchShopNo  string
	)
	for _, target := range uniqueTargets {
		snapshot, innerErr := s.previewRefundTarget(ctx, target)
		if innerErr != nil {
			return nil, innerErr
		}
		if batchOrderNo == "" {
			batchOrderNo = snapshot.GetOrderNo()
			batchShopNo = snapshot.GetShopNo()
		}
		if snapshot.GetOrderNo() != batchOrderNo {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund batch must belong to the same order")
		}
		if snapshot.GetShopNo() != batchShopNo {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund batch must belong to the same shop")
		}
		if innerErr = s.ensureNoActiveRefundCase(ctx, snapshot.GetSubOrderNo()); innerErr != nil {
			return nil, innerErr
		}
		selectedItemNos := normalizeSelectedItemNos(target, snapshot)
		itemNo := pickRefundItemNo(snapshot, target, selectedItemNos)
		afterSaleNo := generateBizNo("AS")
		applyCases = append(applyCases, refundBatchApplyCase{
			afterSaleNo:       afterSaleNo,
			orderNo:           snapshot.GetOrderNo(),
			subOrderNo:        snapshot.GetSubOrderNo(),
			itemNo:            itemNo,
			selectedItemNos:   selectedItemNos,
			snapshot:          snapshot,
			applyRefundAmount: snapshot.GetRefundableAmount(),
		})
	}

	// 先把订单子单冻结到退款中，数据库落单失败时再做补偿恢复，尽量避免卖家发货窗口和退款窗口重叠。
	markedCases := make([]refundBatchApplyCase, 0, len(applyCases))
	for _, applyCase := range applyCases {
		if _, err = s.markOrderRefunding(ctx, applyCase); err != nil {
			s.compensateRefundingMarks(ctx, markedCases)
			return nil, err
		}
		markedCases = append(markedCases, applyCase)
	}

	refundBatchNo := generateBizNo("RB")
	reviewDeadlineAt := gtime.New(time.Now().Add(refundReviewWindow))
	evidenceJSON, _ := json.Marshal(req.GetEvidenceAssetIds())
	now := gtime.Now()
	if err = dao.AfterSaleCase.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, applyCase := range applyCases {
			selectedJSON, _ := json.Marshal(applyCase.selectedItemNos)
			snapshotItem := pickRefundSnapshotItem(applyCase.snapshot, applyCase.itemNo)
			if _, innerErr := tx.Model(dao.AfterSaleCase.Table()).Data(do.AfterSaleCase{
				AfterSaleNo:          applyCase.afterSaleNo,
				OrderNo:              applyCase.orderNo,
				SubOrderNo:           applyCase.subOrderNo,
				ItemNo:               applyCase.itemNo,
				UserId:               userID,
				ShopNo:               applyCase.snapshot.GetShopNo(),
				SpuNo:                snapshotItem.GetSpuNo(),
				SkuNo:                snapshotItem.GetSkuNo(),
				Qty:                  sumRefundSnapshotQty(applyCase.snapshot.GetItems()),
				AfterSaleType:        uint(v1.AfterSaleType_AFTER_SALE_TYPE_REFUND_ONLY),
				AfterSaleStatus:      uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW),
				ApplyRefundAmount:    applyCase.applyRefundAmount,
				ApprovedRefundAmount: 0,
				ReasonCode:           strings.TrimSpace(req.GetReasonCode()),
				ReasonDesc:           strings.TrimSpace(req.GetReasonDesc()),
				EvidenceAssetIdsJson: string(evidenceJSON),
				BuyerRemark:          strings.TrimSpace(req.GetBuyerRemark()),
				RefundBatchNo:        refundBatchNo,
				ScopeCode:            refundScopeSubOrder,
				ReviewDeadlineAt:     reviewDeadlineAt,
				SelectedItemNosJson:  string(selectedJSON),
				PaymentNo:            applyCase.snapshot.GetPaymentNo(),
				Version:              1,
				CreatedAt:            now,
				UpdatedAt:            now,
			}).Insert(); innerErr != nil {
				return gerror.Wrap(innerErr, "create refund after_sale_case failed")
			}
		}
		return s.saveRefundBatchIdempotencyTx(ctx, tx, userID, strings.TrimSpace(req.GetIdempotencyKey()), refundActionApply, refundBatchNo)
	}); err != nil {
		s.compensateRefundingMarks(ctx, markedCases)
		return nil, err
	}

	detail, err := s.buildRefundBatchDetail(ctx, refundBatchNo)
	if err != nil {
		return nil, err
	}
	return &v1.ApplyRefundBatchRes{Detail: detail}, nil
}

// ListMyRefundBatches 按退款批次聚合买家视角的退款列表。
func (s *sAfterSale) ListMyRefundBatches(ctx context.Context, req *v1.ListMyRefundBatchesReq) (*v1.ListMyRefundBatchesRes, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	batches, nextCursor, hasMore, err := s.listRefundBatches(ctx, userID, "", req.GetPageSize(), req.GetNextCursor(), req.GetStatuses())
	if err != nil {
		return nil, err
	}
	return &v1.ListMyRefundBatchesRes{List: batches, NextCursor: nextCursor, HasMore: hasMore}, nil
}

// GetMyRefundBatchDetail 读取买家自己的退款批次详情。
func (s *sAfterSale) GetMyRefundBatchDetail(ctx context.Context, req *v1.GetMyRefundBatchDetailReq) (*v1.GetMyRefundBatchDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetRefundBatchNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_batch_no is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	detail, cases, err := s.buildRefundBatchDetailWithCases(ctx, req.GetRefundBatchNo())
	if err != nil {
		return nil, err
	}
	if len(cases) == 0 || cases[0].UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "refund batch does not belong to user")
	}
	return &v1.GetMyRefundBatchDetailRes{Detail: detail}, nil
}

// CancelRefundBatch 允许买家在卖家还未审核前主动撤回退款，并恢复子单到待发货。
func (s *sAfterSale) CancelRefundBatch(ctx context.Context, req *v1.CancelRefundBatchReq) (*v1.CancelRefundBatchRes, error) {
	if req == nil || strings.TrimSpace(req.GetRefundBatchNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_batch_no is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if hit, batchNo, err := s.replayRefundBatchIdempotency(ctx, userID, strings.TrimSpace(req.GetIdempotencyKey()), refundActionCancel); err != nil {
		return nil, err
	} else if hit {
		detail, innerErr := s.buildRefundBatchDetail(ctx, batchNo)
		if innerErr != nil {
			return nil, innerErr
		}
		return &v1.CancelRefundBatchRes{RefundBatchNo: batchNo, BatchStatus: detail.GetBatch().GetBatchStatus()}, nil
	}

	_, cases, err := s.buildRefundBatchDetailWithCases(ctx, req.GetRefundBatchNo())
	if err != nil {
		return nil, err
	}
	if len(cases) == 0 || cases[0].UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "refund batch does not belong to user")
	}
	if err = ensureRefundBatchAllInStatus(cases, v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW); err != nil {
		return nil, err
	}
	for _, row := range cases {
		if _, err = s.restoreOrderAfterRefundReject(ctx, row.OrderNo, row.SubOrderNo, row.AfterSaleNo); err != nil {
			return nil, err
		}
	}
	err = dao.AfterSaleCase.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, row := range cases {
			if _, innerErr := tx.Model(dao.AfterSaleCase.Table()).
				Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).
				Data(do.AfterSaleCase{
					AfterSaleStatus:  uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_CANCELED),
					CancelReasonCode: "BUYER_CANCELED",
					ClosedAt:         gtime.Now(),
					Version:          row.Version + 1,
				}).Update(); innerErr != nil {
				return gerror.Wrap(innerErr, "cancel refund batch case failed")
			}
		}
		return s.saveRefundBatchIdempotencyTx(ctx, tx, userID, strings.TrimSpace(req.GetIdempotencyKey()), refundActionCancel, strings.TrimSpace(req.GetRefundBatchNo()))
	})
	if err != nil {
		return nil, err
	}
	detail, err := s.buildRefundBatchDetail(ctx, req.GetRefundBatchNo())
	if err != nil {
		return nil, err
	}
	return &v1.CancelRefundBatchRes{
		RefundBatchNo: req.GetRefundBatchNo(),
		BatchStatus:   detail.GetBatch().GetBatchStatus(),
	}, nil
}

// ListShopRefundBatches 返回卖家店铺视角的退款批次列表。
func (s *sAfterSale) ListShopRefundBatches(ctx context.Context, req *v1.ListShopRefundBatchesReq) (*v1.ListShopRefundBatchesRes, error) {
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	batches, nextCursor, hasMore, err := s.listRefundBatches(ctx, 0, req.GetShopNo(), req.GetPageSize(), req.GetNextCursor(), req.GetStatuses())
	if err != nil {
		return nil, err
	}
	return &v1.ListShopRefundBatchesRes{List: batches, NextCursor: nextCursor, HasMore: hasMore}, nil
}

// GetShopRefundBatchDetail 返回卖家店铺视角的退款批次详情。
func (s *sAfterSale) GetShopRefundBatchDetail(ctx context.Context, req *v1.GetShopRefundBatchDetailReq) (*v1.GetShopRefundBatchDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetRefundBatchNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_batch_no is required")
	}
	detail, cases, err := s.buildRefundBatchDetailWithCases(ctx, req.GetRefundBatchNo())
	if err != nil {
		return nil, err
	}
	if len(cases) == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "refund batch not found")
	}
	if shopNo := strings.TrimSpace(req.GetShopNo()); shopNo != "" && cases[0].ShopNo != shopNo {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "refund batch does not belong to shop")
	}
	return &v1.GetShopRefundBatchDetailRes{Detail: detail}, nil
}

// ApproveRefundBatch 卖家同意退款后，串支付退款和订单收口，把批次真正推进到已退款。
func (s *sAfterSale) ApproveRefundBatch(ctx context.Context, req *v1.ApproveRefundBatchReq) (*v1.ApproveRefundBatchRes, error) {
	if req == nil || strings.TrimSpace(req.GetRefundBatchNo()) == "" || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_batch_no/shop_no are required")
	}
	if hit, batchNo, err := s.replayRefundBatchIdempotency(ctx, 0, strings.TrimSpace(req.GetIdempotencyKey()), refundActionApprove); err != nil {
		return nil, err
	} else if hit {
		detail, innerErr := s.buildRefundBatchDetail(ctx, batchNo)
		if innerErr != nil {
			return nil, innerErr
		}
		return &v1.ApproveRefundBatchRes{Detail: detail}, nil
	}

	_, cases, err := s.buildRefundBatchDetailWithCases(ctx, req.GetRefundBatchNo())
	if err != nil {
		return nil, err
	}
	if len(cases) == 0 || cases[0].ShopNo != strings.TrimSpace(req.GetShopNo()) {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "refund batch does not belong to shop")
	}
	if err = ensureRefundBatchAllInStatus(cases, v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW); err != nil {
		return nil, err
	}

	createdTasks := make(map[string]*paymentv1.CreateRefundTaskRes, len(cases))
	for _, row := range cases {
		createRes, innerErr := s.createGatewayRefundTask(ctx, row, fmt.Sprintf("approve_%s", row.AfterSaleNo))
		if innerErr != nil {
			return nil, innerErr
		}
		createdTasks[row.AfterSaleNo] = createRes
	}

	err = dao.AfterSaleCase.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var autoApprovedAt *gtime.Time
		if strings.HasPrefix(strings.TrimSpace(req.GetIdempotencyKey()), "auto_approve_") {
			autoApprovedAt = gtime.Now()
		}
		for _, row := range cases {
			if _, innerErr := tx.Model(dao.AfterSaleCase.Table()).
				Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).
				Data(do.AfterSaleCase{
					AfterSaleStatus:      uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_WAIT_REFUND_TASK),
					ApprovedRefundAmount: row.ApplyRefundAmount,
					SellerReply:          strings.TrimSpace(req.GetSellerReply()),
					AutoApprovedAt:       autoApprovedAt,
					Version:              row.Version + 1,
				}).Update(); innerErr != nil {
				return gerror.Wrap(innerErr, "approve refund batch case failed")
			}
			if innerErr := s.upsertLocalRefundTaskTx(ctx, tx, row, createdTasks[row.AfterSaleNo].GetRefundTaskNo()); innerErr != nil {
				return innerErr
			}
		}
		return s.saveRefundBatchIdempotencyTx(ctx, tx, 0, strings.TrimSpace(req.GetIdempotencyKey()), refundActionApprove, strings.TrimSpace(req.GetRefundBatchNo()))
	})
	if err != nil {
		return nil, err
	}

	for _, row := range cases {
		if _, err = s.ExecuteRefundTask(ctx, &v1.ExecuteRefundTaskReq{RefundTaskNo: createdTasks[row.AfterSaleNo].GetRefundTaskNo()}); err != nil {
			return nil, err
		}
	}
	detail, err := s.buildRefundBatchDetail(ctx, req.GetRefundBatchNo())
	if err != nil {
		return nil, err
	}
	return &v1.ApproveRefundBatchRes{Detail: detail}, nil
}

// RejectRefundBatch 卖家驳回退款时恢复订单履约态，并把批次记为驳回关闭。
func (s *sAfterSale) RejectRefundBatch(ctx context.Context, req *v1.RejectRefundBatchReq) (*v1.RejectRefundBatchRes, error) {
	if req == nil || strings.TrimSpace(req.GetRefundBatchNo()) == "" || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_batch_no/shop_no are required")
	}
	if hit, batchNo, err := s.replayRefundBatchIdempotency(ctx, 0, strings.TrimSpace(req.GetIdempotencyKey()), refundActionReject); err != nil {
		return nil, err
	} else if hit {
		detail, innerErr := s.buildRefundBatchDetail(ctx, batchNo)
		if innerErr != nil {
			return nil, innerErr
		}
		return &v1.RejectRefundBatchRes{Detail: detail}, nil
	}

	_, cases, err := s.buildRefundBatchDetailWithCases(ctx, req.GetRefundBatchNo())
	if err != nil {
		return nil, err
	}
	if len(cases) == 0 || cases[0].ShopNo != strings.TrimSpace(req.GetShopNo()) {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "refund batch does not belong to shop")
	}
	if err = ensureRefundBatchAllInStatus(cases, v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW); err != nil {
		return nil, err
	}
	for _, row := range cases {
		if _, err = s.restoreOrderAfterRefundReject(ctx, row.OrderNo, row.SubOrderNo, row.AfterSaleNo); err != nil {
			return nil, err
		}
	}
	err = dao.AfterSaleCase.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, row := range cases {
			if _, innerErr := tx.Model(dao.AfterSaleCase.Table()).
				Where(dao.AfterSaleCase.Columns().AfterSaleNo, row.AfterSaleNo).
				Data(do.AfterSaleCase{
					AfterSaleStatus:  uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED),
					SellerReply:      strings.TrimSpace(req.GetSellerReply()),
					RejectReasonCode: uint(req.GetRejectReasonCode()),
					ClosedAt:         gtime.Now(),
					Version:          row.Version + 1,
				}).Update(); innerErr != nil {
				return gerror.Wrap(innerErr, "reject refund batch case failed")
			}
		}
		return s.saveRefundBatchIdempotencyTx(ctx, tx, 0, strings.TrimSpace(req.GetIdempotencyKey()), refundActionReject, strings.TrimSpace(req.GetRefundBatchNo()))
	})
	if err != nil {
		return nil, err
	}
	detail, err := s.buildRefundBatchDetail(ctx, req.GetRefundBatchNo())
	if err != nil {
		return nil, err
	}
	return &v1.RejectRefundBatchRes{Detail: detail}, nil
}

// deriveRefundBatchStatus 根据批次内多个子单售后状态推导批次聚合状态。
func deriveRefundBatchStatus(statuses []v1.AfterSaleStatus) v1.AfterSaleStatus {
	if len(statuses) == 0 {
		return v1.AfterSaleStatus_AFTER_SALE_STATUS_UNSPECIFIED
	}
	var (
		hasPending        bool
		hasWaitRefundTask bool
		hasRefunding      bool
		allRefunded       = true
		allRejected       = true
		allCanceled       = true
	)
	for _, status := range statuses {
		if status == v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW {
			hasPending = true
		}
		if status == v1.AfterSaleStatus_AFTER_SALE_STATUS_WAIT_REFUND_TASK {
			hasWaitRefundTask = true
		}
		if status == v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUND_PROCESSING {
			hasRefunding = true
		}
		if status != v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED {
			allRefunded = false
		}
		if status != v1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED {
			allRejected = false
		}
		if status != v1.AfterSaleStatus_AFTER_SALE_STATUS_CANCELED {
			allCanceled = false
		}
	}
	if hasPending {
		return v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW
	}
	if hasWaitRefundTask {
		return v1.AfterSaleStatus_AFTER_SALE_STATUS_WAIT_REFUND_TASK
	}
	if hasRefunding {
		return v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUND_PROCESSING
	}
	if allRefunded {
		return v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED
	}
	if allRejected {
		return v1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED
	}
	if allCanceled {
		return v1.AfterSaleStatus_AFTER_SALE_STATUS_CANCELED
	}
	return v1.AfterSaleStatus_AFTER_SALE_STATUS_CLOSED
}

func (s *sAfterSale) listRefundBatches(ctx context.Context, userID uint64, shopNo string, pageSize int32, nextCursor string, statuses []v1.AfterSaleStatus) ([]*v1.RefundBatch, string, bool, error) {
	pageSize = int32(normalizePageSize(pageSize))
	cursorID, err := parseCursor(nextCursor)
	if err != nil {
		return nil, "", false, err
	}
	cols := dao.AfterSaleCase.Columns()
	model := dao.AfterSaleCase.Ctx(ctx).
		Fields("MAX(id) AS cursor_id, refund_batch_no").
		Where(fmt.Sprintf("%s <> ''", cols.RefundBatchNo)).
		Group(cols.RefundBatchNo).
		OrderDesc("cursor_id").
		Limit(int(pageSize) + 1)
	if userID > 0 {
		model = model.Where(cols.UserId, userID)
	}
	if strings.TrimSpace(shopNo) != "" {
		model = model.Where(cols.ShopNo, shopNo)
	}
	if cursorID > 0 {
		model = model.WhereLT(cols.Id, cursorID)
	}
	if len(statuses) > 0 {
		statusInts := make([]int, 0, len(statuses))
		for _, status := range statuses {
			statusInts = append(statusInts, int(status))
		}
		model = model.WhereIn(cols.AfterSaleStatus, statusInts)
	}
	var cursorRows []refundBatchCursorRow
	if err = model.Scan(&cursorRows); err != nil {
		return nil, "", false, gerror.Wrap(err, "list refund batch cursor rows failed")
	}
	hasMore := false
	if len(cursorRows) > int(pageSize) {
		hasMore = true
		cursorRows = cursorRows[:pageSize]
	}
	list := make([]*v1.RefundBatch, 0, len(cursorRows))
	for _, row := range cursorRows {
		detail, innerErr := s.buildRefundBatchDetail(ctx, row.RefundBatchNo)
		if innerErr != nil {
			return nil, "", false, innerErr
		}
		list = append(list, detail.GetBatch())
	}
	next := ""
	if hasMore && len(cursorRows) > 0 {
		next = fmt.Sprintf("%d", cursorRows[len(cursorRows)-1].CursorID)
	}
	return list, next, hasMore, nil
}

func (s *sAfterSale) buildRefundBatchDetail(ctx context.Context, refundBatchNo string) (*v1.RefundBatchDetail, error) {
	detail, _, err := s.buildRefundBatchDetailWithCases(ctx, refundBatchNo)
	return detail, err
}

func (s *sAfterSale) buildRefundBatchDetailWithCases(ctx context.Context, refundBatchNo string) (*v1.RefundBatchDetail, []*entity.AfterSaleCase, error) {
	cases, err := s.listRefundBatchCases(ctx, refundBatchNo)
	if err != nil {
		return nil, nil, err
	}
	if len(cases) == 0 {
		return nil, nil, gerror.NewCode(gcode.CodeNotFound, "refund batch not found")
	}
	tasks, err := s.listRefundTasksByAfterSaleNos(ctx, collectAfterSaleNos(cases))
	if err != nil {
		return nil, nil, err
	}
	return &v1.RefundBatchDetail{
		Batch:       buildRefundBatch(cases),
		Cases:       toProtoAfterSaleCases(cases),
		RefundTasks: toProtoRefundTasks(tasks),
	}, cases, nil
}

func (s *sAfterSale) listRefundBatchCases(ctx context.Context, refundBatchNo string) ([]*entity.AfterSaleCase, error) {
	var rows []entity.AfterSaleCase
	if err := dao.AfterSaleCase.Ctx(ctx).
		Where(dao.AfterSaleCase.Columns().RefundBatchNo, strings.TrimSpace(refundBatchNo)).
		OrderAsc(dao.AfterSaleCase.Columns().Id).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query refund batch cases failed")
	}
	list := make([]*entity.AfterSaleCase, 0, len(rows))
	for i := range rows {
		list = append(list, &rows[i])
	}
	return list, nil
}

func (s *sAfterSale) listRefundTasksByAfterSaleNos(ctx context.Context, afterSaleNos []string) ([]*entity.RefundTask, error) {
	if len(afterSaleNos) == 0 {
		return nil, nil
	}
	var rows []entity.RefundTask
	if err := dao.RefundTask.Ctx(ctx).
		WhereIn(dao.RefundTask.Columns().AfterSaleNo, afterSaleNos).
		OrderAsc(dao.RefundTask.Columns().Id).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query refund tasks by after_sale_nos failed")
	}
	list := make([]*entity.RefundTask, 0, len(rows))
	for i := range rows {
		list = append(list, &rows[i])
	}
	return list, nil
}

func buildRefundBatch(cases []*entity.AfterSaleCase) *v1.RefundBatch {
	if len(cases) == 0 {
		return nil
	}
	first := cases[0]
	statuses := make([]v1.AfterSaleStatus, 0, len(cases))
	subOrderSeen := make(map[string]struct{}, len(cases))
	subOrderNos := make([]string, 0, len(cases))
	var (
		applyRefundAmount    uint64
		approvedRefundAmount uint64
	)
	for _, row := range cases {
		statuses = append(statuses, v1.AfterSaleStatus(row.AfterSaleStatus))
		applyRefundAmount += row.ApplyRefundAmount
		approvedRefundAmount += row.ApprovedRefundAmount
		if _, ok := subOrderSeen[row.SubOrderNo]; !ok {
			subOrderSeen[row.SubOrderNo] = struct{}{}
			subOrderNos = append(subOrderNos, row.SubOrderNo)
		}
	}
	return &v1.RefundBatch{
		RefundBatchNo:        first.RefundBatchNo,
		OrderNo:              first.OrderNo,
		UserId:               first.UserId,
		ShopNo:               first.ShopNo,
		BatchStatus:          deriveRefundBatchStatus(statuses),
		ApplyRefundAmount:    applyRefundAmount,
		ApprovedRefundAmount: approvedRefundAmount,
		CaseCount:            uint32(len(cases)),
		SubOrderNos:          subOrderNos,
		ReviewDeadlineAt:     toProtoTs(first.ReviewDeadlineAt),
		AutoApprovedAt:       toProtoTs(first.AutoApprovedAt),
		CreatedAt:            toProtoTs(first.CreatedAt),
		UpdatedAt:            toProtoTs(cases[len(cases)-1].UpdatedAt),
	}
}

func toProtoAfterSaleCases(rows []*entity.AfterSaleCase) []*v1.AfterSaleCase {
	list := make([]*v1.AfterSaleCase, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoAfterSaleCase(row))
	}
	return list
}

func toProtoRefundTasks(rows []*entity.RefundTask) []*v1.RefundTask {
	list := make([]*v1.RefundTask, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoRefundTask(row))
	}
	return list
}

func (s *sAfterSale) previewRefundTarget(ctx context.Context, target *v1.ApplyRefundTarget) (*orderv1.RefundSubOrderSnapshot, error) {
	conn, client, err := s.newOrderClient(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	res, err := client.PreviewSubOrderRefund(ctx, &orderv1.PreviewSubOrderRefundReq{
		OrderNo:    strings.TrimSpace(target.GetOrderNo()),
		SubOrderNo: strings.TrimSpace(target.GetSubOrderNo()),
		ItemNo:     strings.TrimSpace(target.GetItemNo()),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "preview sub order refund failed")
	}
	if res == nil || res.GetSnapshot() == nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "refund snapshot is empty")
	}
	return res.GetSnapshot(), nil
}

func (s *sAfterSale) markOrderRefunding(ctx context.Context, applyCase refundBatchApplyCase) (*orderv1.MarkSubOrderRefundingRes, error) {
	conn, client, err := s.newOrderClient(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return client.MarkSubOrderRefunding(ctx, &orderv1.MarkSubOrderRefundingReq{
		OrderNo:    applyCase.orderNo,
		SubOrderNo: applyCase.subOrderNo,
		ItemNo:     applyCase.itemNo,
		RefundNo:   applyCase.afterSaleNo,
	})
}

func (s *sAfterSale) restoreOrderAfterRefundReject(ctx context.Context, orderNo, subOrderNo, refundNo string) (*orderv1.RejectSubOrderRefundRes, error) {
	conn, client, err := s.newOrderClient(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	res, err := client.RejectSubOrderRefund(ctx, &orderv1.RejectSubOrderRefundReq{
		OrderNo:    strings.TrimSpace(orderNo),
		SubOrderNo: strings.TrimSpace(subOrderNo),
		RefundNo:   strings.TrimSpace(refundNo),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "reject sub order refund failed")
	}
	return res, nil
}

func (s *sAfterSale) createGatewayRefundTask(ctx context.Context, row *entity.AfterSaleCase, reason string) (*paymentv1.CreateRefundTaskRes, error) {
	conn, client, err := s.newPaymentClient(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	res, err := client.CreateRefundTask(ctx, &paymentv1.CreateRefundTaskReq{
		AfterSaleNo:    row.AfterSaleNo,
		OrderNo:        row.OrderNo,
		PaymentNo:      row.PaymentNo,
		RefundAmount:   row.ApplyRefundAmount,
		ReasonCode:     row.ReasonCode,
		IdempotencyKey: fmt.Sprintf("%s_%s", reason, row.AfterSaleNo),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "create payment refund task failed")
	}
	return res, nil
}

func (s *sAfterSale) newOrderClient(ctx context.Context) (*grpc.ClientConn, orderv1.InternalOrderServiceClient, error) {
	address := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.orderGrpc", "127.0.0.1:9008").String())
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, gerror.Wrap(err, "dial order grpc failed")
	}
	return conn, orderv1.NewInternalOrderServiceClient(conn), nil
}

func (s *sAfterSale) newPaymentClient(ctx context.Context) (*grpc.ClientConn, paymentv1.InternalPaymentServiceClient, error) {
	address := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.paymentGrpc", "127.0.0.1:9014").String())
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, gerror.Wrap(err, "dial payment grpc failed")
	}
	return conn, paymentv1.NewInternalPaymentServiceClient(conn), nil
}

func (s *sAfterSale) ensureNoActiveRefundCase(ctx context.Context, subOrderNo string) error {
	count, err := dao.AfterSaleCase.Ctx(ctx).
		Where(dao.AfterSaleCase.Columns().SubOrderNo, strings.TrimSpace(subOrderNo)).
		WhereIn(dao.AfterSaleCase.Columns().AfterSaleStatus, []uint{
			uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW),
			uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_WAIT_REFUND_TASK),
			uint(v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUND_PROCESSING),
		}).
		Count()
	if err != nil {
		return gerror.Wrap(err, "query active refund case failed")
	}
	if count > 0 {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "sub order already has processing refund case")
	}
	return nil
}

func (s *sAfterSale) compensateRefundingMarks(ctx context.Context, applyCases []refundBatchApplyCase) {
	for _, applyCase := range applyCases {
		if _, err := s.restoreOrderAfterRefundReject(ctx, applyCase.orderNo, applyCase.subOrderNo, applyCase.afterSaleNo); err != nil {
			g.Log().Warningf(ctx, "[aftersale-svc] rollback refunding mark failed, after_sale_no=%s, err=%+v", applyCase.afterSaleNo, err)
		}
	}
}

func normalizeSelectedItemNos(target *v1.ApplyRefundTarget, snapshot *orderv1.RefundSubOrderSnapshot) []string {
	seen := make(map[string]struct{})
	list := make([]string, 0, len(target.GetSelectedItemNos())+1)
	for _, itemNo := range target.GetSelectedItemNos() {
		itemNo = strings.TrimSpace(itemNo)
		if itemNo == "" {
			continue
		}
		if _, ok := seen[itemNo]; ok {
			continue
		}
		seen[itemNo] = struct{}{}
		list = append(list, itemNo)
	}
	if itemNo := strings.TrimSpace(target.GetItemNo()); itemNo != "" {
		if _, ok := seen[itemNo]; !ok {
			seen[itemNo] = struct{}{}
			list = append(list, itemNo)
		}
	}
	if len(list) == 0 {
		for _, item := range snapshot.GetItems() {
			if strings.TrimSpace(item.GetItemNo()) == "" {
				continue
			}
			list = append(list, item.GetItemNo())
			break
		}
	}
	return list
}

func pickRefundItemNo(snapshot *orderv1.RefundSubOrderSnapshot, target *v1.ApplyRefundTarget, selectedItemNos []string) string {
	if itemNo := strings.TrimSpace(target.GetItemNo()); itemNo != "" {
		return itemNo
	}
	if len(selectedItemNos) > 0 {
		return selectedItemNos[0]
	}
	for _, item := range snapshot.GetItems() {
		if strings.TrimSpace(item.GetItemNo()) != "" {
			return item.GetItemNo()
		}
	}
	return ""
}

func pickRefundSnapshotItem(snapshot *orderv1.RefundSubOrderSnapshot, itemNo string) *orderv1.OrderItemSnapshot {
	for _, item := range snapshot.GetItems() {
		if item.GetItemNo() == itemNo {
			return item
		}
	}
	if len(snapshot.GetItems()) > 0 {
		return snapshot.GetItems()[0]
	}
	return &orderv1.OrderItemSnapshot{}
}

func sumRefundSnapshotQty(items []*orderv1.OrderItemSnapshot) uint {
	var total uint
	for _, item := range items {
		total += uint(item.GetQty())
	}
	if total == 0 {
		total = 1
	}
	return total
}

func collectAfterSaleNos(rows []*entity.AfterSaleCase) []string {
	list := make([]string, 0, len(rows))
	for _, row := range rows {
		list = append(list, row.AfterSaleNo)
	}
	return list
}

func ensureRefundBatchAllInStatus(rows []*entity.AfterSaleCase, status v1.AfterSaleStatus) error {
	for _, row := range rows {
		if v1.AfterSaleStatus(row.AfterSaleStatus) != status {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "refund batch status changed")
		}
	}
	return nil
}
