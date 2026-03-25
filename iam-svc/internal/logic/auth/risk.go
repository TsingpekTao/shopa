package auth

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"google.golang.org/grpc/metadata"
)

// extractRiskMeta 从上下文提取安全元信息。
func extractRiskMeta(ctx context.Context, risk *v1.RiskContext) riskMeta {
	var out riskMeta

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		out.ClientIP = firstMD(md, "x-client-ip")
		out.UserAgent = firstMD(md, "x-user-agent")
		out.RequestID = firstMD(md, "x-request-id")
		out.ForwardedFor = firstMD(md, "x-forwarded-for")
	}
	if out.ClientIP == "" && out.ForwardedFor != "" {
		out.ClientIP = strings.TrimSpace(strings.Split(out.ForwardedFor, ",")[0])
	}

	if risk != nil {
		out.Fingerprint = strings.TrimSpace(risk.GetFingerprint())
		out.DeviceID = strings.TrimSpace(risk.GetDeviceId())
		out.CaptchaToken = strings.TrimSpace(risk.GetCaptchaToken())
	}
	return out
}

// firstMD 读取 metadata 某个 key 的第一个值。
func firstMD(md metadata.MD, key string) string {
	values := md.Get(strings.ToLower(key))
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

// looksLikeAutomation 根据 UA 关键字做轻量自动化工具识别。
func looksLikeAutomation(ua string) bool {
	lower := strings.ToLower(strings.TrimSpace(ua))
	if lower == "" {
		return false
	}
	keywords := []string{"headless", "selenium", "playwright", "puppeteer", "webdriver", "bot"}
	for _, k := range keywords {
		if strings.Contains(lower, k) {
			return true
		}
	}
	return false
}

// bearerFromAuthHeader 从 Authorization 头中提取 bearer token。
func bearerFromAuthHeader(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	const bearer = "bearer "
	if strings.HasPrefix(strings.ToLower(v), bearer) {
		return strings.TrimSpace(v[len(bearer):])
	}
	return v
}

// extractAccessTokenFromMetadata 从上下文提取安全元信息。
func extractAccessTokenFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	auth := firstMD(md, "authorization")
	if auth != "" {
		return bearerFromAuthHeader(auth)
	}
	return firstMD(md, "x-access-token")
}

// jitter 生成随机抖动，避免固定延迟带来的可预测性。
func jitter(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	v := binary.LittleEndian.Uint64(b[:])
	return time.Duration(v % uint64(max))
}

// sleepOnFailure 对失败请求做指数退避延迟，抑制撞库/暴力破解。
func sleepOnFailure(failCount int64) {
	if failCount <= 0 {
		return
	}
	// failCount 越高延迟越长，上限 defaultLoginDelayMax，避免被利用成 DoS。
	delay := defaultLoginDelayBase * time.Duration(1<<minInt64(failCount-1, 4))
	if delay > defaultLoginDelayMax {
		delay = defaultLoginDelayMax
	}
	delay += jitter(120 * time.Millisecond)
	time.Sleep(delay)
}

// minInt64 返回两个 int64 的较小值。
func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
