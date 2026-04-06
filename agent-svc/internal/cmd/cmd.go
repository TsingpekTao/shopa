package cmd

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"github.com/TsingpekTao/shopa/agent-svc/internal/controller/agent"
	"github.com/TsingpekTao/shopa/agent-svc/internal/controller/api"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start agent grpc/http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			go func() {
				cfg := grpcx.Server.NewConfig()
				cfg.Options = append(cfg.Options, []grpc.ServerOption{
					grpcx.Server.ChainUnary(
						grpcx.Server.UnaryValidate,
					),
				}...)
				s := grpcx.Server.New(cfg)
				api.Register(s)
				s.Run()
			}()

			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					agent.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
