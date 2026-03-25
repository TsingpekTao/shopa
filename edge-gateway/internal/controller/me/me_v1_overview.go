package me

import (
	"context"
	"strings"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) GetOverview(ctx context.Context, req *v1.GetOverviewReq) (*v1.GetOverviewRes, error) {
	return service.Bff().BuildMyOverview(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)))
}

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
