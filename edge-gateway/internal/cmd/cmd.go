package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/TsingpekTao/shopa/edge-gateway/internal/controller/admin"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/controller/hello"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/controller/me"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/controller/seller"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/middleware"
	i18nMiddleware "github.com/TsingpekTao/shopa/edge-gateway/internal/middleware/i18n"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// Use global middlewares so proxy can intercept requests
			// even when edge-gateway has no explicit local route binding.
			s.Use(
				middleware.CORS,
				middleware.Proxy,
				i18nMiddleware.Middleware,
				ghttp.MiddlewareHandlerResponse,
			)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Bind(
					admin.NewV1(),
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
