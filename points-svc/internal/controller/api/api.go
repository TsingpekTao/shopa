package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/points-svc/api/v1"
	"github.com/TsingpekTao/shopa/points-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedPointsServiceServer
	points service.IPoints
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{points: service.Points()}
	v1.RegisterPointsServiceServer(s.Server, ctrl)
}

func (c *Controller) InitPointsAccountIfAbsent(ctx context.Context, req *v1.InitPointsAccountIfAbsentReq) (*v1.InitPointsAccountIfAbsentRes, error) {
	return c.points.InitPointsAccountIfAbsent(ctx, req)
}

func (c *Controller) GetPointsByUserId(ctx context.Context, req *v1.GetPointsByUserIdReq) (*v1.GetPointsByUserIdRes, error) {
	return c.points.GetPointsByUserId(ctx, req)
}
