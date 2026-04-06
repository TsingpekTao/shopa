package points

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/points-svc/api/points/v1"
	pb "github.com/TsingpekTao/shopa/points-svc/api/v1"
)

func (c *ControllerV1) InitPointsAccountIfAbsent(ctx context.Context, req *httpv1.InitPointsAccountIfAbsentReq) (*httpv1.InitPointsAccountIfAbsentRes, error) {
	return c.points.InitPointsAccountIfAbsent(ctx, &req.InitPointsAccountIfAbsentReq)
}

func (c *ControllerV1) GetPointsByUserId(ctx context.Context, req *httpv1.GetPointsByUserIdReq) (*httpv1.GetPointsByUserIdRes, error) {
	return c.points.GetPointsByUserId(ctx, &pb.GetPointsByUserIdReq{UserId: req.UserId})
}

func (c *ControllerV1) PreviewOrder(ctx context.Context, req *httpv1.PreviewOrderReq) (*httpv1.PreviewOrderRes, error) {
	pbReq := &pb.PreviewOrderPointsDeductionReq{
		UserId:          req.OrderDraft.UserID,
		OrderNo:         req.OrderDraft.OrderNo,
		OrderAmountCent: int64(req.OrderDraft.PayableAmount),
		RequestedPoints: req.IntentPoints,
		Items:           toPreviewItems(req.OrderDraft.SubOrders),
	}
	pbRes, err := c.points.PreviewOrderPointsDeduction(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	return &httpv1.PreviewOrderRes{
		UsePoints:              pbRes.CanUsePoints,
		PointsUsed:             pbRes.PreviewPoints,
		PointsDiscountAmount:   uint64(maxNonNegative(pbRes.CashDiscountAmountCent)),
		PointsRuleSnapshotJSON: pbRes.RuleSnapshotJson,
		PointsRuleDigest:       pbRes.DeductionDigest,
		SubAllocations:         fromPBSubAllocations(pbRes.SubAllocations),
	}, nil
}

func (c *ControllerV1) LockOrder(ctx context.Context, req *httpv1.LockOrderReq) (*httpv1.LockOrderRes, error) {
	previewRes, err := c.PreviewOrder(ctx, &httpv1.PreviewOrderReq{
		OrderDraft:    req.OrderDraft,
		UsePoints:     req.UsePoints,
		IntentPoints:  req.IntentPoints,
		RequestSource: req.RequestSource,
	})
	if err != nil {
		return nil, err
	}
	pbReq := &pb.LockPointsForOrderReq{
		UserId:                 req.OrderDraft.UserID,
		OrderNo:                req.OrderDraft.OrderNo,
		RequestedPoints:        req.IntentPoints,
		ExpectedCashAmountCent: int64(req.ExpectedPointsCashAmount),
		DeductionDigest:        firstNonEmptyString(req.ExpectedRuleSnapshotDigest, previewRes.PointsRuleDigest),
		RuleSnapshotJson:       previewRes.PointsRuleSnapshotJSON,
		SubAllocations:         toPBSubAllocationsForLock(previewRes.SubAllocations, req.OrderDraft.SubOrders),
		IdempotencyKey:         req.IdempotencyKey,
	}
	pbRes, err := c.points.LockPointsForOrder(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	return &httpv1.LockOrderRes{
		ReservationNo:            pbRes.GetReservation().GetReservationNo(),
		PointsUsed:               pbRes.GetReservation().GetLockedPoints(),
		PointsDiscountAmount:     uint64(maxNonNegative(pbRes.GetReservation().GetLockedCashAmountCent())),
		PointsRuleSnapshotJSON:   pbRes.GetReservation().GetRuleSnapshotJson(),
		PointsRuleSnapshotDigest: pbRes.GetReservation().GetDeductionDigest(),
		SubAllocations:           previewRes.SubAllocations,
	}, nil
}

func (c *ControllerV1) CancelOrder(ctx context.Context, req *httpv1.CancelOrderReq) (*httpv1.SimpleAckRes, error) {
	_, err := c.points.CancelLockedPoints(ctx, &pb.CancelLockedPointsReq{
		ReservationNo:  req.ReservationNo,
		OrderNo:        req.OrderNo,
		ReasonCode:     req.ReasonCode,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.SimpleAckRes{Success: true}, nil
}

func (c *ControllerV1) ConfirmOrder(ctx context.Context, req *httpv1.ConfirmOrderReq) (*httpv1.SimpleAckRes, error) {
	_, err := c.points.ConfirmLockedPoints(ctx, &pb.ConfirmLockedPointsReq{
		ReservationNo:  req.ReservationNo,
		OrderNo:        req.OrderNo,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.SimpleAckRes{Success: true}, nil
}

func (c *ControllerV1) GrantOrder(ctx context.Context, req *httpv1.GrantOrderReq) (*httpv1.SimpleAckRes, error) {
	_, err := c.points.GrantPointsForCompletedOrder(ctx, &pb.GrantPointsForCompletedOrderReq{
		UserId:  req.UserID,
		OrderNo: req.OrderNo,
		Items: []*pb.GrantPointsItem{
			{
				SubOrderNo:         "__ORDER__",
				ItemPaidAmountCent: int64(req.PaidAmount),
			},
		},
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.SimpleAckRes{Success: true}, nil
}

func (c *ControllerV1) ReturnRefund(ctx context.Context, req *httpv1.ReturnRefundReq) (*httpv1.ReturnRefundRes, error) {
	pbRes, err := c.points.ReturnPointsByRefund(ctx, &pb.ReturnPointsByRefundReq{
		UserId:   req.UserID,
		RefundNo: req.RefundNo,
		OrderNo:  req.OrderNo,
		ShopNo:   req.ShopNo,
		Items: []*pb.ReturnPointsItem{
			{
				SubOrderNo:     req.SubOrderNo,
				OrderNo:        req.OrderNo,
				PointsToReturn: req.PointsToReturn,
				CashAmountCent: int64(req.CashAmountCent),
			},
		},
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.ReturnRefundRes{
		Success:            true,
		PointsReturnAmount: pbRes.ReturnedPoints,
	}, nil
}

func (c *ControllerV1) ReverseRefund(ctx context.Context, req *httpv1.ReverseRefundReq) (*httpv1.ReverseRefundRes, error) {
	pbRes, err := c.points.ReverseGrantedPointsByRefund(ctx, &pb.ReverseGrantedPointsByRefundReq{
		UserId:   req.UserID,
		RefundNo: req.RefundNo,
		OrderNo:  req.OrderNo,
		ShopNo:   req.ShopNo,
		Items: []*pb.ReverseGrantedPointsItem{
			{
				SubOrderNo:    req.SubOrderNo,
				ReversePoints: req.PointsToReverse,
			},
		},
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}
	return &httpv1.ReverseRefundRes{
		Success:                true,
		PointsReverseAmount:    pbRes.EffectiveReversedPoints,
		PointsCashOffsetAmount: uint64(maxNonNegative(pbRes.CashOffsetAmountCent)),
		AccountDebtAfter:       pbRes.AccountDebtAfter,
	}, nil
}

func toPreviewItems(subs []httpv1.PointsOrderDraftSub) []*pb.PreviewOrderItem {
	items := make([]*pb.PreviewOrderItem, 0, len(subs))
	for _, sub := range subs {
		items = append(items, &pb.PreviewOrderItem{
			SubOrderNo:     sub.ShopNo,
			ItemAmountCent: int64(sub.PayableAmount),
		})
	}
	return items
}

func fromPBSubAllocations(items []*pb.PointsSubAllocation) []httpv1.PointsSubAllocation {
	result := make([]httpv1.PointsSubAllocation, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		result = append(result, httpv1.PointsSubAllocation{
			ShopNo:               item.SubOrderNo,
			PointsUsed:           item.Points,
			PointsDiscountAmount: uint64(maxNonNegative(item.CashAmountCent)),
		})
	}
	return result
}

func toPBSubAllocationsForLock(items []httpv1.PointsSubAllocation, subs []httpv1.PointsOrderDraftSub) []*pb.PointsSubAllocation {
	amountMap := make(map[string]int64, len(subs))
	for _, sub := range subs {
		amountMap[sub.ShopNo] = int64(sub.PayableAmount)
	}
	result := make([]*pb.PointsSubAllocation, 0, len(items))
	for _, item := range items {
		result = append(result, &pb.PointsSubAllocation{
			SubOrderNo:     item.ShopNo,
			ItemAmountCent: amountMap[item.ShopNo],
			Points:         item.PointsUsed,
			CashAmountCent: int64(item.PointsDiscountAmount),
		})
	}
	return result
}

func maxNonNegative(value int64) int64 {
	if value < 0 {
		return 0
	}
	return value
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
