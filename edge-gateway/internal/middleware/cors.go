package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

func CORS(r *ghttp.Request) {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		origin = "*"
	}

	r.Response.Header().Set("Access-Control-Allow-Origin", origin)
	r.Response.Header().Set("Vary", "Origin")
	r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Request-Id,X-User-Id,X-Account-Status-Code")
	r.Response.Header().Set("Access-Control-Expose-Headers", "X-Request-Id")

	if strings.EqualFold(r.Method, "OPTIONS") {
		r.Response.WriteStatus(204)
		r.ExitAll()
		return
	}

	r.Middleware.Next()
}
