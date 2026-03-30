package points

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/points-svc/api/v1"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedPointsServiceServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterPointsServiceServer(s.Server, &Controller{})
}

func (*Controller) InitPointsAccountIfAbsent(ctx context.Context, req *v1.InitPointsAccountIfAbsentReq) (res *v1.InitPointsAccountIfAbsentRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetPointsByUserId(ctx context.Context, req *v1.GetPointsByUserIdReq) (res *v1.GetPointsByUserIdRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
