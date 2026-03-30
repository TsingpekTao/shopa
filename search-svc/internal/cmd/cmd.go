package cmd

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"github.com/TsingpekTao/shopa/search-svc/internal/controller/api"
	"github.com/TsingpekTao/shopa/search-svc/internal/controller/hello"
)

var (
	// Main 定义搜索服务启动命令。
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start search grpc/http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 异步启动 gRPC 服务，承载内部和公开检索接口。
			go func() {
				// 构造 gRPC 服务器配置对象。
				cfg := grpcx.Server.NewConfig()
				// 注册基础校验中间件链。
				cfg.Options = append(cfg.Options, []grpc.ServerOption{
					grpcx.Server.ChainUnary(
						grpcx.Server.UnaryValidate,
					),
				}...)
				// 创建并启动 gRPC 服务器。
				s := grpcx.Server.New(cfg)
				// 注册业务 gRPC controller。
				api.Register(s)
				// 阻塞运行 gRPC 服务。
				s.Run()
			}()

			// 创建 HTTP 服务器用于健康检查和基础接口。
			s := g.Server()
			// 绑定 HTTP 路由组。
			s.Group("/", func(group *ghttp.RouterGroup) {
				// 启用统一响应中间件。
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				// 绑定默认 hello 接口。
				group.Bind(
					hello.NewV1(),
				)
			})
			// 启动 HTTP 服务。
			s.Run()
			// 正常退出。
			return nil
		},
	}
)
