package points

import (
	v1 "github.com/TsingpekTao/shopa/points-svc/api/v1"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedPointsServiceServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterPointsServiceServer(s.Server, &Controller{})
}
