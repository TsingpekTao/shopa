package middleware

import (
	"net/http"

	"github.com/TsingpekTao/shopa/edge-gateway/internal/consts"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Proxy(r *ghttp.Request) {
	handled, err := service.Proxy().HandleProxyRequest(r.GetCtx(), r)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case gerror.HasCode(err, consts.CodeForbidden):
			status = http.StatusForbidden
		case gerror.HasCode(err, gcode.CodeNotAuthorized):
			status = http.StatusUnauthorized
		case gerror.HasCode(err, gcode.CodeInvalidParameter):
			status = http.StatusBadRequest
		case gerror.HasCode(err, gcode.CodeNotFound):
			status = http.StatusNotFound
		}
		r.Response.WriteStatusExit(status, g.Map{
			"code":    gerror.Code(err).Code(),
			"message": err.Error(),
		})
		return
	}
	if handled {
		r.ExitAll()
		return
	}
	r.Middleware.Next()
}
