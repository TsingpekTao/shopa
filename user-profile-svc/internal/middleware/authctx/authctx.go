package authctx

import (
	"net/http"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/ctxkey"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// UserIDFromHeader 从请求头读取用户 ID 并注入上下文。
func UserIDFromHeader(r *ghttp.Request) {
	// 先从常见请求头读取 user id。
	// 生产环境应由 iam-svc 的 JWT 校验直接提供 user id。
	var (
		raw = r.GetHeader("X-User-Id")
	)
	if raw == "" {
		raw = r.GetHeader("User-Id")
	}
	if raw == "" {
		raw = r.GetHeader("Uid")
	}
	if raw == "" {
		raw = r.GetHeader("X-Uid")
	}

	// 若缺少用户标识，则拒绝请求，防止匿名访问 me/*。
	if raw == "" {
		// 向客户端写入未授权状态码。
		r.Response.WriteStatus(http.StatusUnauthorized)
		// 设置 error 以便统一响应与日志记录。
		r.SetError(gerror.NewCode(gcode.CodeNotAuthorized, "missing X-User-Id header"))
		// 中止请求流程。
		r.Exit()
	}

	// 转换为 uint64 并校验合法值。
	userID := gconv.Uint64(raw)
	if userID == 0 {
		// 通过 gerror 便于框架记录堆栈。
		r.SetError(gerror.NewCode(gcode.CodeInvalidParameter, "invalid X-User-Id header"))
		r.Response.WriteStatus(http.StatusBadRequest)
		// 中止请求流程。
		r.Exit()
	}

	// 将用户 ID 写入请求上下文，供后续中间件/处理器使用。
	r.SetCtxVar(ctxkey.UserIDKey{}, userID)

	// 继续执行后续中间件与处理器。
	r.Middleware.Next()
}
