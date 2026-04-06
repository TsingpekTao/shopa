package v1

import (
	pb "github.com/TsingpekTao/shopa/promotion-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListAvailableCouponsReq struct {
	g.Meta      `path:"/v1/promotions/coupons/available" method:"get" tags:"Promotion" summary:"List available coupons"`
	UserId      uint64 `json:"userId" in:"query"`
	OrderAmount uint64 `json:"orderAmount" in:"query"`
}

type ListAvailableCouponsRes = pb.ListAvailableCouponsRes

type CalculateOrderDiscountReq struct {
	g.Meta `path:"/v1/promotions/discounts/calculate" method:"post" tags:"Promotion" summary:"Calculate order discount"`
	pb.CalculateOrderDiscountReq
}

type CalculateOrderDiscountRes = pb.CalculateOrderDiscountRes

type LockCouponForOrderReq struct {
	g.Meta         `path:"/v1/internal/promotions/coupon-locks" method:"post" tags:"PromotionInternal" summary:"Lock coupon for order"`
	UserId         uint64                 `json:"userId"`
	OrderNo        string                 `json:"orderNo"`
	CouponNos      []string               `json:"couponNos"`
	LockExpireAt   *timestamppb.Timestamp `json:"lockExpireAt"`
	IdempotencyKey string                 `json:"idempotencyKey"`
}

type LockCouponForOrderRes = pb.LockCouponForOrderRes

type ConfirmCouponUsageReq struct {
	g.Meta `path:"/v1/internal/promotions/coupon-locks/confirm" method:"post" tags:"PromotionInternal" summary:"Confirm coupon usage"`
	pb.ConfirmCouponUsageReq
}

type ConfirmCouponUsageRes = pb.ConfirmCouponUsageRes

type ReleaseCouponLockReq struct {
	g.Meta `path:"/v1/internal/promotions/coupon-locks/release" method:"post" tags:"PromotionInternal" summary:"Release coupon lock"`
	pb.ReleaseCouponLockReq
}

type ReleaseCouponLockRes = pb.ReleaseCouponLockRes

type CreateCampaignReq struct {
	g.Meta `path:"/v1/admin/promotions/campaigns" method:"post" tags:"PromotionAdmin" summary:"Create promotion campaign"`
	pb.CreateCampaignReq
}

type CreateCampaignRes = pb.CreateCampaignRes
