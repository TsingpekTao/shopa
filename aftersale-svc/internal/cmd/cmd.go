package cmd

import (
	"context"

	"github.com/TsingpekTao/shopa/aftersale-svc/internal/controller/aftersale"
	_ "github.com/TsingpekTao/shopa/aftersale-svc/internal/logic"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/worker"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"github.com/TsingpekTao/shopa/aftersale-svc/internal/controller/api"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/controller/hello"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start aftersale gRPC and HTTP servers",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			worker.StartRefundAutoApproveWorker(ctx)

			go func() {
				c := grpcx.Server.NewConfig()
				c.Options = append(c.Options, []grpc.ServerOption{
					grpcx.Server.ChainUnary(
						grpcx.Server.UnaryValidate,
					),
				}...)
				s := grpcx.Server.New(c)
				api.Register(s)
				s.Run()
			}()

			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
					aftersale.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
