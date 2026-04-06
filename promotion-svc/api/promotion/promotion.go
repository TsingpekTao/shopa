package promotion

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/promotion-svc/api/promotion/v1"
)

type IPromotionV1 interface {
	ListAvailableCoupons(ctx context.Context, req *v1.ListAvailableCouponsReq) (res *v1.ListAvailableCouponsRes, err error)
	CalculateOrderDiscount(ctx context.Context, req *v1.CalculateOrderDiscountReq) (res *v1.CalculateOrderDiscountRes, err error)
	LockCouponForOrder(ctx context.Context, req *v1.LockCouponForOrderReq) (res *v1.LockCouponForOrderRes, err error)
	ConfirmCouponUsage(ctx context.Context, req *v1.ConfirmCouponUsageReq) (res *v1.ConfirmCouponUsageRes, err error)
	ReleaseCouponLock(ctx context.Context, req *v1.ReleaseCouponLockReq) (res *v1.ReleaseCouponLockRes, err error)
	CreateCampaign(ctx context.Context, req *v1.CreateCampaignReq) (res *v1.CreateCampaignRes, err error)
}
