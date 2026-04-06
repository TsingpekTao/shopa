package agent

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc/metadata"
)

func withRequestMetadata(ctx context.Context) context.Context {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ctx
	}
	md := metadata.Pairs(
		"x-user-id", req.GetHeader("X-User-Id"),
		"x-operator-user-id", req.GetHeader("X-Operator-User-Id"),
		"x-request-id", req.GetHeader("X-Request-Id"),
		"x-idempotency-key", req.GetHeader("X-Idempotency-Key"),
		"x-forwarded-for", req.GetHeader("X-Forwarded-For"),
		"x-client-ip", req.GetClientIp(),
		"x-user-agent", req.GetHeader("User-Agent"),
		"authorization", req.GetHeader("Authorization"),
		"x-access-token", req.GetHeader("X-Access-Token"),
	)
	return metadata.NewIncomingContext(ctx, md)
}
