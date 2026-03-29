package admin

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/admin/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) GetOverview(ctx context.Context, req *v1.GetOverviewReq) (*v1.GetOverviewRes, error) {
	return service.Bff().BuildAdminOverview(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)))
}

func (c *ControllerV1) GetDashboardOverview(ctx context.Context, req *v1.GetDashboardOverviewReq) (*v1.GetDashboardOverviewRes, error) {
	return service.Bff().BuildAdminDashboardOverview(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)))
}

func (c *ControllerV1) GetShopInsights(ctx context.Context, req *v1.GetShopInsightsReq) (*v1.GetShopInsightsRes, error) {
	return service.Bff().BuildAdminShopInsights(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), req.ShopNo)
}

func (c *ControllerV1) ListConversations(ctx context.Context, req *v1.ListConversationsReq) (*v1.ListConversationsRes, error) {
	return service.Bff().BuildAdminConversationList(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)))
}
