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

// riskMeta 汇总一次请求的风控上下文，供登录/注册/审计日志复用。
type riskMeta struct {
	ClientIP     string
	UserAgent    string
	RequestID    string
	ForwardedFor string
	Fingerprint  string
	DeviceID     string
	CaptchaToken string
}

// mfaChallengePayload 是 MFA challenge 在缓存中的序列化载荷。
// 仅保存完成二次校验所需最小字段，避免缓存泄露时暴露过多信息。
type mfaChallengePayload struct {
	UserID     uint64 `json:"user_id"`
	Phone      string `json:"phone"`
	Identifier string `json:"identifier"`
}

// registerEventPayload 是用户注册事件写入 outbox 的消息体。
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

// now 统一返回 UTC 当前时间，避免时区混用导致 token/过期判断偏差。
func now() time.Time {
	return time.Now().UTC()
}

// newEventID 生成全局唯一事件 ID。
func newEventID() string {
	return uuid.NewString()
}

// sceneKey 把短信场景枚举映射为稳定字符串键，用于缓存 key 拼接。
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

// normalizePhone 对手机号做最小归一化，当前策略仅 trim 空白。
func normalizePhone(phone string) string {
	return strings.TrimSpace(phone)
}

// normalizeIdentifier 归一化登录标识（手机号或邮箱）。
// 邮箱场景统一转小写，减少同一账号的等价输入差异。
func normalizeIdentifier(identifier string) string {
	return strings.ToLower(strings.TrimSpace(identifier))
}

// toProtoTimestamp 把 gtime 转为 protobuf timestamp。
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

// randomDigits 生成固定长度数字串（短信验证码）。
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

// randomToken 生成无填充 Base32 随机串，默认大写便于日志和人工输入识别。
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

// makeInitDisplayName 生成初始展示名。
// 组合规则：`usernamePrefix + 随机后缀`，通过随机后缀降低重名概率。
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
