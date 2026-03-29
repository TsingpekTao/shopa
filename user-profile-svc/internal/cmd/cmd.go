package cmd

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/controller/api"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/middleware/i18n"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/router"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/worker"
)

var (
	// Main 是服务启动主命令，统一拉起 gRPC 与 HTTP。
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start grpc server and http server",
		Func:  mainFunc,
	}
)

// mainFunc 启动消费者、gRPC 服务与 HTTP 服务。
func mainFunc(ctx context.Context, parser *gcmd.Parser) (err error) {
	// 启动“注册后初始化资料”消息消费者。
	worker.StartRegisterInitConsumer(ctx)

	// gRPC 使用单独协程启动，避免阻塞当前主协程，确保 HTTP 也能同时启动。
	go func() {
		// 创建 gRPC 服务器配置，用于挂载拦截器链。
		c := grpcx.Server.NewConfig()
		// 在 Unary 拦截器链中接入参数校验，尽早拦截非法请求。
		c.Options = append(c.Options, []grpc.ServerOption{
			grpcx.Server.ChainUnary(
				grpcx.Server.UnaryValidate,
			),
		}...)
		// 基于配置构造 gRPC 服务实例。
		s := grpcx.Server.New(c)
		// 注册 user-profile 的 RPC 处理器。
		api.Register(s)
		// 启动并阻塞在 gRPC 服务协程中。
		s.Run()
	}()

	// 构建 HTTP 服务（Swagger/OpenAPI + 对外 HTTP 接口）。
	httpServer := g.Server()
	// 使用统一响应中间件，保证返回结构一致。
	httpServer.Use(i18n.Middleware)
	httpServer.Use(ghttp.MiddlewareHandlerResponse)
	// 按模块注册 HTTP 路由分组。
	router.RegisterHTTP(httpServer)
	// 自定义 OpenAPI 文档信息。
	router.EnhanceOpenAPIDoc(httpServer)
	// 启动 HTTP 并阻塞主协程。
	httpServer.Run()
	return nil
}
