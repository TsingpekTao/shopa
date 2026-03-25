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

// extractRiskMeta 组装风控上下文。
// 数据来源：
// 1) gRPC metadata（网关透传头，如 client-ip/ua/request-id）。
// 2) 业务请求体内的 RiskContext（设备指纹、验证码 token 等）。
// 该函数只负责提取和归一化，不承载策略判定逻辑。
func extractRiskMeta(ctx context.Context, risk *v1.RiskContext) riskMeta {
	var out riskMeta

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		out.ClientIP = firstMD(md, "x-client-ip")
		out.UserAgent = firstMD(md, "x-user-agent")
		out.RequestID = firstMD(md, "x-request-id")
		out.ForwardedFor = firstMD(md, "x-forwarded-for")
	}
	// 代理链路下，如果 x-client-ip 缺失，则退化为 X-Forwarded-For 的首个地址。
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

// firstMD 读取 metadata 中指定 key 的第一个值。
// gRPC metadata 的 key 不区分大小写，这里统一转小写读取。
func firstMD(md metadata.MD, key string) string {
	values := md.Get(strings.ToLower(key))
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

// looksLikeAutomation 基于 UA 关键字做轻量自动化工具识别。
// 这里只做低成本初筛，命中后通常要求补充验证码，避免误杀正常用户。
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

// bearerFromAuthHeader 从 Authorization 头提取 Bearer token。
// 兼容两种输入：
// 说明：1) "Bearer <token>"
// 2) 直接传入 token（无前缀）。
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

// extractAccessTokenFromMetadata 从 metadata 提取 access token。
// 优先级：authorization > x-access-token。
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

// jitter 生成 [0, max) 随机抖动，用于打散重试/延迟时间，避免雪崩。
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

// sleepOnFailure 在登录失败后执行指数退避。
// 安全意图：增加暴力破解成本；系统意图：通过抖动避免请求在固定时间点齐发。
func sleepOnFailure(failCount int64) {
	if failCount <= 0 {
		return
	}
	// 最大只放大到 2^4，防止延迟无限增长导致可用性问题。
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
