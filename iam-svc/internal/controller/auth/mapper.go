package auth

import (
	"context"
	"time"

	authv1 "github.com/TsingpekTao/shopa/iam-svc/api/auth/v1"
	iamv1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// withRequestMetadata 从 HTTP 请求中提取 metadata 以透传到 gRPC。
func withRequestMetadata(ctx context.Context) context.Context {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ctx
	}
	md := metadata.Pairs(
		"x-client-ip", req.GetClientIp(),
		"x-user-agent", req.GetHeader("User-Agent"),
		"x-request-id", req.GetHeader("X-Request-Id"),
		"x-forwarded-for", req.GetHeader("X-Forwarded-For"),
		"authorization", req.GetHeader("Authorization"),
		"x-access-token", req.GetHeader("X-Access-Token"),
	)
	return metadata.NewIncomingContext(ctx, md)
}

// toProtoRisk 把 HTTP 风控结构映射为 protobuf 风控结构。
func toProtoRisk(in *authv1.RiskContext) *iamv1.RiskContext {
	if in == nil {
		return nil
	}
	return &iamv1.RiskContext{
		Fingerprint:  in.Fingerprint,
		DeviceId:     in.DeviceID,
		CaptchaToken: in.CaptchaToken,
	}
}

// toHTTPTokenPair 将 protobuf TokenPair 转换为 HTTP 侧结构。
func toHTTPTokenPair(in *iamv1.TokenPair) *authv1.TokenPair {
	if in == nil {
		return nil
	}
	return &authv1.TokenPair{
		TokenType:        in.TokenType,
		AccessToken:      in.AccessToken,
		AccessExpiresIn:  in.AccessExpiresIn,
		RefreshToken:     in.RefreshToken,
		RefreshExpiresIn: in.RefreshExpiresIn,
		Sid:              in.Sid,
	}
}

// toHTTPAuthResult 将 gRPC AuthResult 映射为 HTTP 友好的结构。
func toHTTPAuthResult(in *iamv1.AuthResult) *authv1.AuthResult {
	if in == nil {
		return nil
	}
	out := &authv1.AuthResult{}
	if token := in.GetTokenPair(); token != nil {
		out.TokenPair = toHTTPTokenPair(token)
	}
	if challenge := in.GetMfaChallenge(); challenge != nil {
		out.MfaChallenge = &authv1.MfaChallenge{
			ChallengeID: challenge.ChallengeId,
			ExpireAt:    toTimeString(challenge.ExpireAt),
		}
	}
	return out
}

// toHTTPSession 将 protobuf SessionSummary 转为 HTTP 会话摘要。
func toHTTPSession(in *iamv1.SessionSummary) *authv1.SessionSummary {
	if in == nil {
		return nil
	}
	roles := make([]*authv1.RoleItem, 0, len(in.Roles))
	for _, role := range in.Roles {
		roles = append(roles, &authv1.RoleItem{
			RoleCode:  int32(role.RoleCode),
			ScopeType: int32(role.ScopeType),
			ScopeID:   role.ScopeId,
		})
	}

	var membership *authv1.MembershipSummary
	if in.Membership != nil {
		membership = &authv1.MembershipSummary{
			LevelCode: in.Membership.LevelCode,
			Points:    in.Membership.Points,
			ExpireAt:  toTimeString(in.Membership.ExpireAt),
		}
	}

	return &authv1.SessionSummary{
		UserID:        in.UserId,
		AccountStatus: int32(in.AccountStatus),
		Roles:         roles,
		Membership:    membership,
		LastLoginAt:   toTimeString(in.LastLoginAt),
		LastLoginIP:   in.LastLoginIp,
	}
}

// toTimeString 把 protobuf Timestamp 格式化为 RFC3339 字符串。
func toTimeString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().UTC().Format(time.RFC3339)
}
