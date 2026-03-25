package seller

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
)

type ISellerV1 interface {
	GetWorkbench(ctx context.Context, req *v1.GetWorkbenchReq) (res *v1.GetWorkbenchRes, err error)
	GetShopDashboard(ctx context.Context, req *v1.GetShopDashboardReq) (res *v1.GetShopDashboardRes, err error)
}
