package seller

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
)

type ISellerV1 interface {
	GetWorkbench(ctx context.Context, req *v1.GetWorkbenchReq) (res *v1.GetWorkbenchRes, err error)
	GetShopDashboard(ctx context.Context, req *v1.GetShopDashboardReq) (res *v1.GetShopDashboardRes, err error)
	GetSalesAnalytics(ctx context.Context, req *v1.GetSalesAnalyticsReq) (res *v1.GetSalesAnalyticsRes, err error)
	ListRefundBatches(ctx context.Context, req *v1.ListRefundBatchesReq) (res *v1.ListRefundBatchesRes, err error)
	GetRefundBatchDetail(ctx context.Context, req *v1.GetRefundBatchDetailReq) (res *v1.GetRefundBatchDetailRes, err error)
	ApproveRefundBatch(ctx context.Context, req *v1.ApproveRefundBatchReq) (res *v1.ApproveRefundBatchRes, err error)
	RejectRefundBatch(ctx context.Context, req *v1.RejectRefundBatchReq) (res *v1.RejectRefundBatchRes, err error)
}
