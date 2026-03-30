package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/promotion-svc/api/v1"
	"github.com/TsingpekTao/shopa/promotion-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedBuyerPromotionServiceServer
	v1.UnimplementedInternalPromotionServiceServer
	v1.UnimplementedAdminPromotionServiceServer
}

func Register(s *grpcx.GrpcServer) {
	controller := &Controller{}
	v1.RegisterBuyerPromotionServiceServer(s.Server, controller)
	v1.RegisterInternalPromotionServiceServer(s.Server, controller)
	v1.RegisterAdminPromotionServiceServer(s.Server, controller)
}

func (*Controller) ListAvailableCoupons(ctx context.Context, req *v1.ListAvailableCouponsReq) (res *v1.ListAvailableCouponsRes, err error) {
	return service.Promotion().ListAvailableCoupons(ctx, req)
}

func (*Controller) CalculateOrderDiscount(ctx context.Context, req *v1.CalculateOrderDiscountReq) (res *v1.CalculateOrderDiscountRes, err error) {
	return service.Promotion().CalculateOrderDiscount(ctx, req)
}

func (*Controller) LockCouponForOrder(ctx context.Context, req *v1.LockCouponForOrderReq) (res *v1.LockCouponForOrderRes, err error) {
	return service.Promotion().LockCouponForOrder(ctx, req)
}

func (*Controller) ConfirmCouponUsage(ctx context.Context, req *v1.ConfirmCouponUsageReq) (res *v1.ConfirmCouponUsageRes, err error) {
	return service.Promotion().ConfirmCouponUsage(ctx, req)
}

func (*Controller) ReleaseCouponLock(ctx context.Context, req *v1.ReleaseCouponLockReq) (res *v1.ReleaseCouponLockRes, err error) {
	return service.Promotion().ReleaseCouponLock(ctx, req)
}

func (*Controller) CreateCampaign(ctx context.Context, req *v1.CreateCampaignReq) (res *v1.CreateCampaignRes, err error) {
	return service.Promotion().CreateCampaign(ctx, req)
}
