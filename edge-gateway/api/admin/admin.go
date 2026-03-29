package admin

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/admin/v1"
)

type IAdminV1 interface {
	GetOverview(ctx context.Context, req *v1.GetOverviewReq) (res *v1.GetOverviewRes, err error)
	GetDashboardOverview(ctx context.Context, req *v1.GetDashboardOverviewReq) (res *v1.GetDashboardOverviewRes, err error)
	GetShopInsights(ctx context.Context, req *v1.GetShopInsightsReq) (res *v1.GetShopInsightsRes, err error)
	ListConversations(ctx context.Context, req *v1.ListConversationsReq) (res *v1.ListConversationsRes, err error)
}
