package points

import (
	"context"
	"strings"

	httpv1 "github.com/TsingpekTao/shopa/points-svc/api/points/v1"
	pb "github.com/TsingpekTao/shopa/points-svc/api/v1"
	"github.com/TsingpekTao/shopa/points-svc/internal/model/entity"
)

// PreviewOrderPointsDeduction 兼容旧版 PB 接口，复用 HTTP 预览逻辑生成订单积分抵扣结果。
func (s *sPoints) PreviewOrderPointsDeduction(ctx context.Context, req *pb.PreviewOrderPointsDeductionReq) (*pb.PreviewOrderPointsDeductionRes, error) {
	legacyReq := &httpv1.PreviewOrderReq{
		OrderDraft:   buildLegacyOrderDraft(req.GetUserId(), req.GetOrderNo(), req.GetOrderAmountCent(), req.GetItems()),
		UsePoints:    req.GetRequestedPoints() > 0,
		IntentPoints: req.GetRequestedPoints(),
	}
	legacyRes, err := s.previewOrderHTTP(ctx, legacyReq)
	if err != nil {
		return nil, err
	}

	account, err := s.ensureAccount(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.PreviewOrderPointsDeductionRes{
		CanUsePoints:           legacyRes.UsePoints,
		ReasonCode:             buildPreviewReasonCode(req, account, legacyRes),
		AvailablePoints:        toCompatBalance(account),
		MaxUsablePoints:        legacyRes.PointsUsed,
		PreviewPoints:          legacyRes.PointsUsed,
		CashDiscountAmountCent: int64(legacyRes.PointsDiscountAmount),
		DeductionDigest:        legacyRes.PointsRuleDigest,
		RuleSnapshotJson:       legacyRes.PointsRuleSnapshotJSON,
		SubAllocations:         toPBSubAllocations(legacyRes.SubAllocations, req.GetItems()),
	}, nil
}

// LockPointsForOrder 将 PB 协议请求转成 HTTP 领域请求，沿用统一的预留锁定流程。
func (s *sPoints) LockPointsForOrder(ctx context.Context, req *pb.LockPointsForOrderReq) (*pb.LockPointsForOrderRes, error) {
	legacyReq := &httpv1.LockOrderReq{
		OrderDraft:                 buildLegacyOrderDraftFromAllocations(req.GetUserId(), req.GetOrderNo(), req.GetShopNo(), req.GetSubAllocations()),
		UsePoints:                  req.GetRequestedPoints() > 0,
		IntentPoints:               req.GetRequestedPoints(),
		IdempotencyKey:             req.GetIdempotencyKey(),
		ExpectedPointsCashAmount:   uint64(maxInt64(req.GetExpectedCashAmountCent(), 0)),
		ExpectedRuleSnapshotDigest: req.GetDeductionDigest(),
	}
	legacyRes, err := s.lockOrderHTTP(ctx, legacyReq)
	if err != nil {
		return nil, err
	}

	reservation, _ := s.findReservation(ctx, legacyRes.ReservationNo, req.GetOrderNo())
	account, err := s.ensureAccount(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.LockPointsForOrderRes{
		Reservation:           toPBReservation(reservation),
		AccountAvailableAfter: account.AvailableBalance,
		AccountFrozenAfter:    account.FrozenBalance,
	}, nil
}

// ConfirmLockedPoints 确认已锁定积分，并返回确认后的账户与预留快照。
func (s *sPoints) ConfirmLockedPoints(ctx context.Context, req *pb.ConfirmLockedPointsReq) (*pb.ConfirmLockedPointsRes, error) {
	if _, err := s.confirmOrderHTTP(ctx, &httpv1.ConfirmOrderReq{
		ReservationNo:  req.GetReservationNo(),
		OrderNo:        req.GetOrderNo(),
		IdempotencyKey: req.GetIdempotencyKey(),
	}); err != nil {
		return nil, err
	}

	reservation, err := s.findReservation(ctx, req.GetReservationNo(), req.GetOrderNo())
	if err != nil {
		return nil, err
	}
	account, err := s.ensureAccount(ctx, reservation.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.ConfirmLockedPointsRes{
		Reservation:           toPBReservation(reservation),
		AccountAvailableAfter: account.AvailableBalance,
		AccountFrozenAfter:    account.FrozenBalance,
		DebtPointsAfter:       debtFromAvailable(account.AvailableBalance),
	}, nil
}

// CancelLockedPoints 取消已锁定积分，并把释放后的账户状态回填到 PB 响应。
func (s *sPoints) CancelLockedPoints(ctx context.Context, req *pb.CancelLockedPointsReq) (*pb.CancelLockedPointsRes, error) {
	if _, err := s.cancelOrderHTTP(ctx, &httpv1.CancelOrderReq{
		ReservationNo:  req.GetReservationNo(),
		OrderNo:        req.GetOrderNo(),
		ReasonCode:     req.GetReasonCode(),
		IdempotencyKey: req.GetIdempotencyKey(),
	}); err != nil {
		return nil, err
	}

	reservation, err := s.findReservation(ctx, req.GetReservationNo(), req.GetOrderNo())
	if err != nil {
		return nil, err
	}
	account, err := s.ensureAccount(ctx, reservation.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.CancelLockedPointsRes{
		Reservation:           toPBReservation(reservation),
		AccountAvailableAfter: account.AvailableBalance,
		AccountFrozenAfter:    account.FrozenBalance,
	}, nil
}

// GrantPointsForCompletedOrder 聚合订单项实付金额，并复用授予逻辑完成订单返积分。
func (s *sPoints) GrantPointsForCompletedOrder(ctx context.Context, req *pb.GrantPointsForCompletedOrderReq) (*pb.GrantPointsForCompletedOrderRes, error) {
	var paidAmount uint64
	for _, item := range req.GetItems() {
		if item == nil {
			continue
		}
		if item.GetItemPaidAmountCent() > 0 {
			paidAmount += uint64(item.GetItemPaidAmountCent())
		}
	}

	if _, err := s.grantOrderHTTP(ctx, &httpv1.GrantOrderReq{
		OrderNo:        req.GetOrderNo(),
		UserID:         req.GetUserId(),
		PaidAmount:     paidAmount,
		IdempotencyKey: req.GetIdempotencyKey(),
	}); err != nil {
		return nil, err
	}

	account, err := s.ensureAccount(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.GrantPointsForCompletedOrderRes{
		GrantedPoints:         paidAmount,
		AccountAvailableAfter: account.AvailableBalance,
		DebtPointsAfter:       debtFromAvailable(account.AvailableBalance),
	}, nil
}

// ReturnPointsByRefund 按退款明细逐条回退消费积分，兼容旧协议的退款返还入口。
func (s *sPoints) ReturnPointsByRefund(ctx context.Context, req *pb.ReturnPointsByRefundReq) (*pb.ReturnPointsByRefundRes, error) {
	var returned uint64
	for _, item := range req.GetItems() {
		if item == nil {
			continue
		}
		legacyRes, err := s.returnRefundHTTP(ctx, &httpv1.ReturnRefundReq{
			RefundNo:       req.GetRefundNo(),
			OrderNo:        firstNonEmpty(item.GetOrderNo(), req.GetOrderNo()),
			SubOrderNo:     item.GetSubOrderNo(),
			ShopNo:         req.GetShopNo(),
			UserID:         req.GetUserId(),
			PointsToReturn: item.GetPointsToReturn(),
			CashAmountCent: uint64(maxInt64(item.GetCashAmountCent(), 0)),
			IdempotencyKey: req.GetIdempotencyKey(),
		})
		if err != nil {
			return nil, err
		}
		returned += legacyRes.PointsReturnAmount
	}

	account, err := s.ensureAccount(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &pb.ReturnPointsByRefundRes{
		ReturnedPoints:        returned,
		AccountAvailableAfter: account.AvailableBalance,
		DebtPointsAfter:       debtFromAvailable(account.AvailableBalance),
		GraceBucketCreated:    returned > 0 && account.AvailableBalance >= 0,
	}, nil
}

// ReverseGrantedPointsByRefund 按退款项冲回已赠送积分，并返回实际冲回与现金补差信息。
func (s *sPoints) ReverseGrantedPointsByRefund(ctx context.Context, req *pb.ReverseGrantedPointsByRefundReq) (*pb.ReverseGrantedPointsByRefundRes, error) {
	var (
		totalRequested uint64
		totalReversed  uint64
		totalOffset    uint64
		remainCash     uint64
	)
	for _, item := range req.GetItems() {
		if item == nil {
			continue
		}
		totalRequested += item.GetReversePoints()
		legacyRes, err := s.reverseRefundHTTP(ctx, &httpv1.ReverseRefundReq{
			RefundNo:             req.GetRefundNo(),
			OrderNo:              req.GetOrderNo(),
			SubOrderNo:           item.GetSubOrderNo(),
			ShopNo:               req.GetShopNo(),
			UserID:               req.GetUserId(),
			PointsToReverse:      item.GetReversePoints(),
			ApprovedRefundAmount: remainCash,
			IdempotencyKey:       req.GetIdempotencyKey(),
		})
		if err != nil {
			return nil, err
		}
		totalReversed += legacyRes.PointsReverseAmount
		totalOffset += legacyRes.PointsCashOffsetAmount
		if legacyRes.PointsCashOffsetAmount <= remainCash {
			remainCash -= legacyRes.PointsCashOffsetAmount
		} else {
			remainCash = 0
		}
	}

	account, err := s.ensureAccount(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	deficit := uint64(0)
	if totalRequested > totalReversed {
		deficit = totalRequested - totalReversed
	}

	return &pb.ReverseGrantedPointsByRefundRes{
		ReversedPoints:          totalRequested,
		EffectiveReversedPoints: totalReversed,
		DeficitPoints:           deficit,
		CashOffsetAmountCent:    int64(totalOffset),
		AccountAvailableAfter:   account.AvailableBalance,
		AccountDebtAfter:        debtFromAvailable(account.AvailableBalance),
	}, nil
}

// buildLegacyOrderDraft 把旧 PB 订单项聚合成 HTTP 侧使用的订单草稿结构。
func buildLegacyOrderDraft(userID uint64, orderNo string, orderAmountCent int64, items []*pb.PreviewOrderItem) *httpv1.PointsOrderDraft {
	subOrders := make(map[string]uint64)
	total := uint64(maxInt64(orderAmountCent, 0))
	var sum uint64
	for _, item := range items {
		if item == nil {
			continue
		}
		subNo := strings.TrimSpace(item.GetSubOrderNo())
		if subNo == "" {
			subNo = "__DEFAULT__"
		}
		amount := uint64(maxInt64(item.GetItemAmountCent(), 0))
		subOrders[subNo] += amount
		sum += amount
	}
	if total == 0 {
		total = sum
	}
	if len(subOrders) == 0 && total > 0 {
		subOrders["__DEFAULT__"] = total
	}

	result := &httpv1.PointsOrderDraft{
		OrderNo:       orderNo,
		UserID:        userID,
		PayableAmount: total,
		SubOrders:     make([]httpv1.PointsOrderDraftSub, 0, len(subOrders)),
	}
	for subNo, amount := range subOrders {
		result.SubOrders = append(result.SubOrders, httpv1.PointsOrderDraftSub{
			ShopNo:        subNo,
			GoodsAmount:   amount,
			PayableAmount: amount,
		})
	}
	return result
}

// buildLegacyOrderDraftFromAllocations 按子单分摊结果重建订单草稿，供锁定流程复用。
func buildLegacyOrderDraftFromAllocations(userID uint64, orderNo, shopNo string, allocations []*pb.PointsSubAllocation) *httpv1.PointsOrderDraft {
	subOrders := make([]httpv1.PointsOrderDraftSub, 0, len(allocations))
	var total uint64
	for _, item := range allocations {
		if item == nil {
			continue
		}
		subNo := strings.TrimSpace(item.GetSubOrderNo())
		if subNo == "" {
			subNo = firstNonEmpty(strings.TrimSpace(shopNo), "__DEFAULT__")
		}
		amount := uint64(maxInt64(item.GetItemAmountCent(), 0))
		total += amount
		subOrders = append(subOrders, httpv1.PointsOrderDraftSub{
			ShopNo:        subNo,
			GoodsAmount:   amount,
			PayableAmount: amount,
		})
	}
	return &httpv1.PointsOrderDraft{
		OrderNo:       orderNo,
		UserID:        userID,
		PayableAmount: total,
		SubOrders:     subOrders,
	}
}

// toPBSubAllocations 将 HTTP 侧的按店铺分配结果转换回 PB 协议定义。
func toPBSubAllocations(items []httpv1.PointsSubAllocation, previews []*pb.PreviewOrderItem) []*pb.PointsSubAllocation {
	amountMap := make(map[string]int64, len(previews))
	for _, item := range previews {
		if item == nil {
			continue
		}
		amountMap[item.GetSubOrderNo()] += item.GetItemAmountCent()
	}

	result := make([]*pb.PointsSubAllocation, 0, len(items))
	for _, item := range items {
		subNo := item.ShopNo
		result = append(result, &pb.PointsSubAllocation{
			SubOrderNo:     subNo,
			ItemAmountCent: amountMap[subNo],
			Points:         item.PointsUsed,
			CashAmountCent: int64(item.PointsDiscountAmount),
		})
	}
	return result
}

// toPBReservation 将数据库中的预留记录映射为 PB 返回结构。
func toPBReservation(row *entity.PointsReservation) *pb.ReservationSnapshot {
	if row == nil {
		return nil
	}
	return &pb.ReservationSnapshot{
		ReservationNo:         row.ReservationNo,
		UserId:                row.UserId,
		OrderNo:               row.OrderNo,
		ShopNo:                row.ShopNo,
		ReservationStatusCode: row.ReservationStatusCode,
		LockedPoints:          row.LockedPoints,
		LockedCashAmountCent:  row.LockedCashAmountCent,
		DeductionDigest:       row.DeductionDigest,
		RuleSnapshotJson:      row.RuleSnapshotJson,
		IdempotencyKey:        row.IdempotencyKey,
		ExpireAt:              formatTime(row.ExpireAt),
		ConfirmedAt:           formatTime(row.ConfirmedAt),
		CanceledAt:            formatTime(row.CanceledAt),
		CancelReasonCode:      row.CancelReasonCode,
		CreatedAt:             formatTime(row.CreatedAt),
		UpdatedAt:             formatTime(row.UpdatedAt),
	}
}

// buildPreviewReasonCode 在不可抵扣时给出兼容旧客户端的原因码。
func buildPreviewReasonCode(req *pb.PreviewOrderPointsDeductionReq, account *entity.PointsAccount, legacyRes *httpv1.PreviewOrderRes) string {
	if legacyRes != nil && legacyRes.UsePoints {
		return ""
	}
	if req.GetRequestedPoints() == 0 {
		return "POINTS_NOT_REQUESTED"
	}
	if account == nil || account.AvailableBalance <= 0 {
		return "NO_AVAILABLE_POINTS"
	}
	if req.GetOrderAmountCent() <= 0 {
		return "ORDER_AMOUNT_INVALID"
	}
	return "PREVIEW_POINTS_ZERO"
}

// maxInt64 保证金额类字段不会被负数带入后续积分计算。
func maxInt64(v int64, fallback int64) int64 {
	if v < fallback {
		return fallback
	}
	return v
}

// firstNonEmpty 返回第一个非空白字符串，用于兼容缺失字段的旧请求。
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
