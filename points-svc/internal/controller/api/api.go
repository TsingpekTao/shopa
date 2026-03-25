package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/points-svc/api/points/v1"
	"github.com/TsingpekTao/shopa/points-svc/internal/service/points"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedPointsServiceServer
	points *points.Service
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{points: points.New()}
	v1.RegisterPointsServiceServer(s.Server, ctrl)
}

func (c *Controller) InitPointsAccountIfAbsent(ctx context.Context, req *v1.InitPointsAccountIfAbsentReq) (*v1.InitPointsAccountIfAbsentRes, error) {
	return c.points.InitPointsAccountIfAbsent(ctx, req)
}

func (c *Controller) GetPointsByUserId(ctx context.Context, req *v1.GetPointsByUserIdReq) (*v1.GetPointsByUserIdRes, error) {
	return c.points.GetPointsByUserId(ctx, req)
}
