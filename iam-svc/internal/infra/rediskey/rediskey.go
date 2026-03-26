package rediskey

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/genv"
)

const (
	TTLRegisterIPWindow = 2 * time.Hour
	TTLSmsSendLock      = 60 * time.Second
	TTLSmsCode          = 5 * time.Minute
	TTLSmsCodeAttempt   = 5 * time.Minute
	TTLLoginFailCounter = 30 * time.Minute
	TTLLoginLock        = 30 * time.Minute
	TTLMFAChallenge     = 10 * time.Minute
)

// Builder 负责构建 IAM 领域的 Redis Key。
type Builder struct {
	env     string
	project string
	service string
	prefix  string
}

// New 从配置创建默认 Key Builder。
func New() *Builder {
	var (
		ctx     = context.Background()
		env     = strings.TrimSpace(g.Cfg().MustGet(ctx, "app.env", "").String())
		project = strings.TrimSpace(g.Cfg().MustGet(ctx, "app.project", "shopa").String())
		service = strings.TrimSpace(g.Cfg().MustGet(ctx, "app.service", "iam").String())
	)
	if env == "" {
		env = strings.TrimSpace(genv.Get("GF_ENV", "dev").String())
	}
	if project == "" {
		project = "shopa"
	}
	if service == "" {
		service = "iam"
	}
	return NewWith(env, project, service)
}

// NewWith 按传入环境信息创建 Key Builder。
func NewWith(env, project, service string) *Builder {
	env = sanitizePart(env)
	project = sanitizePart(project)
	service = sanitizePart(service)
	prefix := strings.Join([]string{env, project, service}, ":")
	return &Builder{
		env:     env,
		project: project,
		service: service,
		prefix:  prefix,
	}
}

// Prefix 返回统一前缀。
func (b *Builder) Prefix() string {
	return b.prefix
}

// Key 组装完整 Redis Key。
func (b *Builder) Key(parts ...string) string {
	out := make([]string, 0, 1+len(parts))
	out = append(out, b.prefix)
	for _, p := range parts {
		out = append(out, sanitizePart(p))
	}
	return strings.Join(out, ":")
}

// RegIPLimitKey 生成注册 IP 限流键。
func (b *Builder) RegIPLimitKey(ip string, t time.Time) string {
	return b.Key("reg", "ip", ip, t.Format("2006010215"))
}

// SmsLockKey 生成短信发送锁键。
func (b *Builder) SmsLockKey(scene, phone string) string {
	return b.Key("sms", "lock", scene, phone)
}

// SmsCodeKey 生成短信验证码键。
func (b *Builder) SmsCodeKey(scene, phone string) string {
	return b.Key("sms", "code", scene, phone)
}

// SmsCodeAttemptKey 生成短信验证码尝试次数键。
func (b *Builder) SmsCodeAttemptKey(scene, phone string) string {
	return b.Key("sms", "code", "attempt", scene, phone)
}

// LoginFailKey 生成登录失败计数键。
func (b *Builder) LoginFailKey(identifier string) string {
	return b.Key("login", "fail", identifier)
}

// LoginLockKey 生成登录锁定键。
func (b *Builder) LoginLockKey(identifier string) string {
	return b.Key("login", "lock", identifier)
}

// MFAChallengeKey 生成 MFA 挑战键。
func (b *Builder) MFAChallengeKey(challengeID string) string {
	return b.Key("mfa", "challenge", challengeID)
}

// sanitizePart 对 key 片段做安全规范化。
func sanitizePart(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "_"
	}
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "\t", "_")
	s = strings.ReplaceAll(s, "\n", "_")
	return s
}

// MustFormat 按格式拼接字符串。
func MustFormat(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
