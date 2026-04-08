package seller

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) GetWorkbench(ctx context.Context, req *v1.GetWorkbenchReq) (*v1.GetWorkbenchRes, error) {
	return service.Bff().BuildSellerWorkbench(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)))
}

func (c *ControllerV1) GetShopDashboard(ctx context.Context, req *v1.GetShopDashboardReq) (*v1.GetShopDashboardRes, error) {
	return service.Bff().BuildSellerShopDashboard(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), req.ShopNo)
}

func (c *ControllerV1) GetSalesAnalytics(ctx context.Context, req *v1.GetSalesAnalyticsReq) (*v1.GetSalesAnalyticsRes, error) {
	return service.Bff().BuildSellerSalesAnalytics(
		ctx,
		extractAccessTokenFromRequest(g.RequestFromCtx(ctx)),
		req.ShopNo,
		req.Range,
	)
}
