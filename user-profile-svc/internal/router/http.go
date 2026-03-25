package router

import (
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/controller/hello"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/controller/user"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/middleware/authctx"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
)

// RegisterHTTP 澶勭悊娉ㄥ唽涓绘祦绋嬪強鍒濆鍖栧姩浣溿€
func RegisterHTTP(s *ghttp.Server) {
	// Root group for service routes.
	s.Group("/", func(group *ghttp.RouterGroup) {
		// Public endpoints (no auth).
		group.Bind(
			hello.NewV1(),
		)

		// "Me" endpoints (require user id; in real production: iam-svc token -> user id).
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

// EnhanceOpenAPIDoc 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func EnhanceOpenAPIDoc(s *ghttp.Server) {
	// Read openapi generator instance from server.
	openapi := s.GetOpenApi()
	// Align OpenAPI common response with GoFrame default handler response wrapper.
	openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	openapi.Config.CommonResponseDataField = `Data`
	// Basic API info for Swagger UI.
	openapi.Info = goai.Info{
		Title:       `Shopa User Profile Service`,
		Description: `HTTP APIs for user profile/address + auto OpenAPI/Swagger docs.`,
	}
}
