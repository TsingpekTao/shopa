package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/v1"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedMeBffServiceServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterMeBffServiceServer(s.Server, &Controller{})
}

func (*Controller) GetMyOverview(ctx context.Context, req *v1.GetMyOverviewReq) (res *v1.GetMyOverviewRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetSellerWorkbench(ctx context.Context, req *v1.GetSellerWorkbenchReq) (res *v1.GetSellerWorkbenchRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetSellerShopDashboard(ctx context.Context, req *v1.GetSellerShopDashboardReq) (res *v1.GetSellerShopDashboardRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
