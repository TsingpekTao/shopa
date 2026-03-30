// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/promotion-svc/api/v1"
)

type (
	IPromotion interface {
		// ListAvailableCoupons 查询用户可用券列表。
		ListAvailableCoupons(ctx context.Context, req *v1.ListAvailableCouponsReq) (*v1.ListAvailableCouponsRes, error)
		// CalculateOrderDiscount 计算订单优惠金额。
		CalculateOrderDiscount(ctx context.Context, req *v1.CalculateOrderDiscountReq) (*v1.CalculateOrderDiscountRes, error)
		// LockCouponForOrder 锁定订单使用券。
		LockCouponForOrder(ctx context.Context, req *v1.LockCouponForOrderReq) (*v1.LockCouponForOrderRes, error)
		// ConfirmCouponUsage 确认券使用成功。
		ConfirmCouponUsage(ctx context.Context, req *v1.ConfirmCouponUsageReq) (*v1.ConfirmCouponUsageRes, error)
		// ReleaseCouponLock 释放订单券锁。
		ReleaseCouponLock(ctx context.Context, req *v1.ReleaseCouponLockReq) (*v1.ReleaseCouponLockRes, error)
		// CreateCampaign 创建促销活动。
		CreateCampaign(ctx context.Context, req *v1.CreateCampaignReq) (*v1.CreateCampaignRes, error)
	}
)

var (
	localPromotion IPromotion
)

func Promotion() IPromotion {
	if localPromotion == nil {
		panic("implement not found for interface IPromotion, forgot register?")
	}
	return localPromotion
}

func RegisterPromotion(i IPromotion) {
	localPromotion = i
}
