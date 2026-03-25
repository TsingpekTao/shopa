// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	mev1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
	sellerv1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
)

type (
	IBff interface {
		BuildMyOverview(ctx context.Context, accessToken string) (*mev1.GetOverviewRes, error)
		BuildSellerWorkbench(ctx context.Context, accessToken string) (*sellerv1.GetWorkbenchRes, error)
		BuildSellerShopDashboard(ctx context.Context, accessToken string, shopNo string) (*sellerv1.GetShopDashboardRes, error)
	}
)

var (
	localBff IBff
)

func Bff() IBff {
	if localBff == nil {
		panic("implement not found for interface IBff, forgot register?")
	}
	return localBff
}

func RegisterBff(i IBff) {
	localBff = i
}
