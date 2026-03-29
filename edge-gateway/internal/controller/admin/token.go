package admin

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

func extractAccessTokenFromRequest(r *ghttp.Request) string {
	if r == nil {
		return ""
	}
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(auth) >= len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		return strings.TrimSpace(auth[len(prefix):])
	}
	return auth
}
