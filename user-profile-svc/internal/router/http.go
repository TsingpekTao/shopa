package router

import (
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/controller/hello"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/controller/user"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/middleware/authctx"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
)

// RegisterHTTP 注册 HTTP 路由与相关分组。
func RegisterHTTP(s *ghttp.Server) {
	// 按服务模块建立根路由组。
	s.Group("/", func(group *ghttp.RouterGroup) {
		// 注册无需认证的公开接口。
		group.Bind(
			hello.NewV1(),
		)

		// "Me" 接口需要 user id，生产环境应由 iam-svc token 提供。
		group.Group("/v1/me", func(me *ghttp.RouterGroup) {
			me.Middleware(
				authctx.UserIDFromHeader,
			)
			me.Bind(
				user.NewV1(),
			)
		})
	})
}

// EnhanceOpenAPIDoc 配置 OpenAPI 文档与通用响应。
func EnhanceOpenAPIDoc(s *ghttp.Server) {
	// 从服务器读取 OpenAPI 生成器实例。
	openapi := s.GetOpenApi()
	// 同步 OpenAPI 的公共响应结构与 GoFrame 默认响应。
	openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	openapi.Config.CommonResponseDataField = `Data`
	// 设置 Swagger UI 的基本信息。
	openapi.Info = goai.Info{
		Title:       `Shopa User Profile Service`,
		Description: `HTTP APIs for user profile/address + auto OpenAPI/Swagger docs.`,
	}
}
