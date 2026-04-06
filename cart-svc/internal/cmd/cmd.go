package cmd

import (
	"context"
	"net/http"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"github.com/TsingpekTao/shopa/cart-svc/internal/controller/api"
	"github.com/TsingpekTao/shopa/cart-svc/internal/controller/cart"
	"github.com/TsingpekTao/shopa/cart-svc/internal/controller/hello"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start cart gRPC and HTTP servers",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
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
				group.Middleware(func(r *ghttp.Request) {
					if r.Method != http.MethodGet && r.ContentLength != 0 {
						r.MakeBodyRepeatableRead(true)
					}
					r.Middleware.Next()
				})
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
					cart.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
