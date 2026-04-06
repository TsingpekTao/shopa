package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// riskMeta 汇总单次请求的风控上下文，供登录、注册和审计日志复用。
type riskMeta struct {
	ClientIP     string
	UserAgent    string
	RequestID    string
	ForwardedFor string
	Fingerprint  string
	DeviceID     string
	CaptchaToken string
}

// mfaChallengePayload 是 MFA challenge 在缓存中的最小化载荷。
// 这里只保留二次校验所需字段，避免缓存泄露时暴露过多信息。
type mfaChallengePayload struct {
	UserID     uint64 `json:"user_id"`
	Phone      string `json:"phone"`
	Identifier string `json:"identifier"`
}

// registerEventPayload 是注册成功后写入 outbox 的事件内容。
type registerEventPayload struct {
	EventID         string `json:"event_id"`
	EventVersion    string `json:"event_version"`
	UserID          uint64 `json:"user_id"`
	InitDisplayName string `json:"init_display_name"`
	RegisterChannel string `json:"register_channel"`
	OccurredAt      string `json:"occurred_at"`
}

// gctx 返回后台任务可复用的基础 context。
func gctx() context.Context {
	return context.Background()
}

// now 统一返回 UTC 当前时间。
func now() time.Time {
	return time.Now().UTC()
}

// newEventID 生成全局唯一事件 ID。
func newEventID() string {
	return uuid.NewString()
}

// sceneKey 将短信场景枚举映射为稳定字符串键。
func sceneKey(scene v1.SmsScene) string {
	switch scene {
	case v1.SmsScene_SMS_SCENE_REGISTER:
		return "register"
	case v1.SmsScene_SMS_SCENE_LOGIN:
		return "login"
	case v1.SmsScene_SMS_SCENE_MFA:
		return "mfa"
	case v1.SmsScene_SMS_SCENE_RESET_PASSWORD:
		return "reset_password"
	default:
		return "unknown"
	}
}

// normalizePhone 对手机号做最小归一化，当前只去除首尾空白。
func normalizePhone(phone string) string {
	return strings.TrimSpace(phone)
}

// normalizeIdentifier 对登录标识做归一化，便于邮箱登录统一比较。
func normalizeIdentifier(identifier string) string {
	return strings.ToLower(strings.TrimSpace(identifier))
}

// toProtoTimestamp 将 GoFrame 的 gtime 转为 protobuf timestamp。
func toProtoTimestamp(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(time.Unix(t.Timestamp(), 0).UTC())
}

// sha256Hex 计算字符串的 SHA-256 十六进制摘要。
func sha256Hex(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// randomDigits 生成指定长度的纯数字验证码。
func randomDigits(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	var b strings.Builder
	b.Grow(n)
	for _, v := range bytes {
		b.WriteByte(byte('0' + (v % 10)))
	}
	return b.String(), nil
}

// randomToken 生成 Base32 编码的随机字符串，用于 challenge 或 token 辅助字段。
func randomToken(size int) (string, error) {
	if size <= 0 {
		size = 16
	}
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	return strings.ToUpper(enc.EncodeToString(buf)), nil
}

// makeInitDisplayName 生成注册时的默认展示名。
func (s *Service) makeInitDisplayName() (string, error) {
	suffix, err := randomToken(4)
	if err != nil {
		return "", err
	}
	if len(suffix) > 6 {
		suffix = suffix[:6]
	}
	prefix := strings.TrimSpace(s.usernamePrefix)
	if prefix == "" {
		prefix = defaultUsernamePrefix
	}
	return fmt.Sprintf("%s%s", prefix, suffix), nil
}
