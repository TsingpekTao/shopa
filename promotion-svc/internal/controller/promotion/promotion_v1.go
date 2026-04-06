package promotion

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/promotion-svc/api/promotion/v1"
	pb "github.com/TsingpekTao/shopa/promotion-svc/api/v1"
)

func (c *ControllerV1) ListAvailableCoupons(ctx context.Context, req *httpv1.ListAvailableCouponsReq) (*httpv1.ListAvailableCouponsRes, error) {
	return c.promotion.ListAvailableCoupons(ctx, &pb.ListAvailableCouponsReq{
		UserId:      req.UserId,
		OrderAmount: req.OrderAmount,
	})
}

func (c *ControllerV1) CalculateOrderDiscount(ctx context.Context, req *httpv1.CalculateOrderDiscountReq) (*httpv1.CalculateOrderDiscountRes, error) {
	return c.promotion.CalculateOrderDiscount(ctx, &req.CalculateOrderDiscountReq)
}

func (c *ControllerV1) LockCouponForOrder(ctx context.Context, req *httpv1.LockCouponForOrderReq) (*httpv1.LockCouponForOrderRes, error) {
	return c.promotion.LockCouponForOrder(ctx, &pb.LockCouponForOrderReq{
		UserId:         req.UserId,
		OrderNo:        req.OrderNo,
		CouponNos:      req.CouponNos,
		LockExpireAt:   req.LockExpireAt,
		IdempotencyKey: req.IdempotencyKey,
	})
}

func (c *ControllerV1) ConfirmCouponUsage(ctx context.Context, req *httpv1.ConfirmCouponUsageReq) (*httpv1.ConfirmCouponUsageRes, error) {
	return c.promotion.ConfirmCouponUsage(ctx, &req.ConfirmCouponUsageReq)
}

func (c *ControllerV1) ReleaseCouponLock(ctx context.Context, req *httpv1.ReleaseCouponLockReq) (*httpv1.ReleaseCouponLockRes, error) {
	return c.promotion.ReleaseCouponLock(ctx, &req.ReleaseCouponLockReq)
}

func (c *ControllerV1) CreateCampaign(ctx context.Context, req *httpv1.CreateCampaignReq) (*httpv1.CreateCampaignRes, error) {
	return c.promotion.CreateCampaign(ctx, &req.CreateCampaignReq)
}
