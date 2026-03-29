package cmd

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"github.com/TsingpekTao/shopa/iam-svc/internal/controller/api"
	"github.com/TsingpekTao/shopa/iam-svc/internal/controller/auth"
	"github.com/TsingpekTao/shopa/iam-svc/internal/controller/hello"
	"github.com/TsingpekTao/shopa/iam-svc/internal/middleware/i18n"
	"github.com/TsingpekTao/shopa/iam-svc/internal/middleware/response"
	"github.com/TsingpekTao/shopa/iam-svc/internal/service"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start IAM gRPC and HTTP servers",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			service.Auth().StartBackgroundWorkers(ctx)

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
				group.Middleware(i18n.Middleware, response.Middleware)
				group.Bind(
					hello.NewV1(),
					auth.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
