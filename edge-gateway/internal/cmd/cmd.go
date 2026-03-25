package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/TsingpekTao/shopa/edge-gateway/internal/controller/hello"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/controller/me"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/controller/seller"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					middleware.Proxy,
					ghttp.MiddlewareHandlerResponse,
				)
				group.Bind(
					hello.NewV1(),
					me.NewV1(),
					seller.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
