package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func CORS(r *ghttp.Request) {
	origin := strings.TrimSpace(r.Header.Get("Origin"))

	allowedOrigins := g.Cfg().MustGet(r.Context(), "gateway.cors.allowedOrigins", "http://127.0.0.1:3000,http://127.0.0.1:3100,http://127.0.0.1:3200").Strings()
	allowedMethods := g.Cfg().MustGet(r.Context(), "gateway.cors.allowedMethods", "GET,POST,PUT,PATCH,DELETE,OPTIONS").String()
	allowedHeaders := g.Cfg().MustGet(r.Context(), "gateway.cors.allowedHeaders", "Authorization,Content-Type,X-Request-Id,X-User-Id,X-Account-Status-Code,X-Idempotency-Key,X-Shop-No").String()
	allowCredentials := g.Cfg().MustGet(r.Context(), "gateway.cors.allowCredentials", true).Bool()

	isAllowed := isOriginAllowed(origin, allowedOrigins)
	if origin != "" && isAllowed {
		r.Response.Header().Set("Access-Control-Allow-Origin", origin)
	}

	r.Response.Header().Set("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
	if allowCredentials {
		r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	r.Response.Header().Set("Access-Control-Allow-Methods", allowedMethods)
	r.Response.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	r.Response.Header().Set("Access-Control-Expose-Headers", "X-Request-Id")

	if strings.EqualFold(r.Method, "OPTIONS") {
		if origin != "" && !isAllowed {
			r.Response.WriteStatus(403)
			r.ExitAll()
			return
		}
		r.Response.WriteStatus(204)
		r.ExitAll()
		return
	}

	r.Middleware.Next()
}

func isOriginAllowed(origin string, whitelist []string) bool {
	if strings.TrimSpace(origin) == "" {
		return true
	}
	for _, item := range whitelist {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if item == "*" || strings.EqualFold(item, origin) {
			return true
		}
	}
	return false
}
