package cmd

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"github.com/TsingpekTao/shopa/risk-svc/internal/controller/api"
	"github.com/TsingpekTao/shopa/risk-svc/internal/controller/hello"
	"github.com/TsingpekTao/shopa/risk-svc/internal/controller/risk"
)

var (
	// Main 启动 risk-svc 的 HTTP 与 gRPC 服务。
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start risk grpc/http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 启动 gRPC 服务并注册 risk API。
			go func() {
				// 创建 gRPC 配置对象。
				c := grpcx.Server.NewConfig()
				// 加入统一参数校验中间件。
				c.Options = append(c.Options, []grpc.ServerOption{
					grpcx.Server.ChainUnary(grpcx.Server.UnaryValidate),
				}...)
				// 初始化 gRPC 服务实例。
				s := grpcx.Server.New(c)
				// 注册 risk controller。
				api.Register(s)
				// 启动 gRPC 服务。
				s.Run()
			}()

			// 启动 HTTP 服务用于探活与调试。
			s := g.Server()
			// 绑定根路由分组。
			s.Group("/", func(group *ghttp.RouterGroup) {
				// 启用统一响应封装中间件。
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				// 绑定 hello 路由。
				group.Bind(
					hello.NewV1(),
					risk.NewV1(),
				)
			})
			// 启动 HTTP 服务。
			s.Run()
			// 返回空错误表示正常退出。
			return nil
		},
	}
)
