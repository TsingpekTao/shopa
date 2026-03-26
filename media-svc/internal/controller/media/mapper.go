package media

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc/metadata"
)

// withRequestMetadata 将 HTTP 请求头透传为 gRPC metadata，复用 logic 层统一上下文解析逻辑。
func withRequestMetadata(ctx context.Context) context.Context {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ctx
	}
	md := metadata.Pairs(
		"x-user-id", req.GetHeader("X-User-Id"),
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
