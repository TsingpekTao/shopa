// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 说明：如需手动维护此接口文件，可删除这些注释。
// ================================================================================

package service

import (
	"context"

	aftersalev1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	adminv1 "github.com/TsingpekTao/shopa/edge-gateway/api/admin/v1"
	mev1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
	sellerv1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
)

type (
	// IBff 定义网关 BFF 聚合门面：
	// - BuildMyOverview: 用户侧“我的概览”聚合。
	// - BuildSellerWorkbench: 卖家工作台聚合。
	// - BuildSellerShopDashboard: 店铺仪表盘聚合。
	IBff interface {
		BuildMyOverview(ctx context.Context, accessToken string) (*mev1.GetOverviewRes, error)
		BuildAdminOverview(ctx context.Context, accessToken string) (*adminv1.GetOverviewRes, error)
		BuildAdminDashboardOverview(ctx context.Context, accessToken string) (*adminv1.GetDashboardOverviewRes, error)
		BuildAdminShopInsights(ctx context.Context, accessToken string, shopNo string) (*adminv1.GetShopInsightsRes, error)
		BuildAdminConversationList(ctx context.Context, accessToken string) (*adminv1.ListConversationsRes, error)
		BuildSellerWorkbench(ctx context.Context, accessToken string) (*sellerv1.GetWorkbenchRes, error)
		BuildSellerShopDashboard(ctx context.Context, accessToken string, shopNo string) (*sellerv1.GetShopDashboardRes, error)
		BuildSellerSalesAnalytics(ctx context.Context, accessToken string, shopNo string, rangeCode string) (*sellerv1.GetSalesAnalyticsRes, error)
		BuildBuyerRefundPreview(ctx context.Context, accessToken string, req *mev1.PreviewRefundReq) (*mev1.PreviewRefundRes, error)
		ApplyBuyerRefundBatch(ctx context.Context, accessToken string, req *aftersalev1.ApplyRefundBatchReq) (*aftersalev1.ApplyRefundBatchRes, error)
		ListBuyerRefundBatches(ctx context.Context, accessToken string, req *mev1.ListRefundBatchesReq) (*aftersalev1.ListMyRefundBatchesRes, error)
		GetBuyerRefundBatchDetail(ctx context.Context, accessToken string, req *mev1.GetRefundBatchDetailReq) (*aftersalev1.GetMyRefundBatchDetailRes, error)
		CancelBuyerRefundBatch(ctx context.Context, accessToken string, req *aftersalev1.CancelRefundBatchReq) (*aftersalev1.CancelRefundBatchRes, error)
		ListSellerRefundBatches(ctx context.Context, accessToken string, req *sellerv1.ListRefundBatchesReq) (*aftersalev1.ListShopRefundBatchesRes, error)
		GetSellerRefundBatchDetail(ctx context.Context, accessToken string, req *sellerv1.GetRefundBatchDetailReq) (*aftersalev1.GetShopRefundBatchDetailRes, error)
		ApproveSellerRefundBatch(ctx context.Context, accessToken string, req *aftersalev1.ApproveRefundBatchReq) (*aftersalev1.ApproveRefundBatchRes, error)
		RejectSellerRefundBatch(ctx context.Context, accessToken string, req *aftersalev1.RejectRefundBatchReq) (*aftersalev1.RejectRefundBatchRes, error)
	}
)

var (
	// localBff 保存已注册的 BFF 实现。
	localBff IBff
)

// Bff 返回 BFF 服务实现；未注册时 panic，避免请求落到空实现。
func Bff() IBff {
	if localBff == nil {
		panic("implement not found for interface IBff, forgot register?")
	}
	return localBff
}

// RegisterBff 注册 BFF 实现。
func RegisterBff(i IBff) {
	localBff = i
}
