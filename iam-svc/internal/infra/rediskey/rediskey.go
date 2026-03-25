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

// Builder 閹稿绮烘稉鈧崜宥囩磻鐟欏嫬鍨弸鍕紦 Redis Key閿涘矂浼╅崗宥咁樋閻滎垰顣?婢舵碍婀囬崝鈥冲暱缁愪降鈧
type Builder struct {
	env     string
	project string
	service string
	prefix  string
}

// New 鏋勯€犲嚱鏁帮細鍒涘缓骞惰繑鍥炴湇鍔″疄渚嬨€
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

// NewWith 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
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

// Prefix 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (b *Builder) Prefix() string {
	return b.prefix
}

// Key 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (b *Builder) Key(parts ...string) string {
	out := make([]string, 0, 1+len(parts))
	out = append(out, b.prefix)
	for _, p := range parts {
		out = append(out, sanitizePart(p))
	}
	return strings.Join(out, ":")
}

// RegIPLimitKey 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (b *Builder) RegIPLimitKey(ip string, t time.Time) string {
	return b.Key("reg", "ip", ip, t.Format("2006010215"))
}

// SmsLockKey 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (b *Builder) SmsLockKey(scene, phone string) string {
	return b.Key("sms", "lock", scene, phone)
}

// SmsCodeKey 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (b *Builder) SmsCodeKey(scene, phone string) string {
	return b.Key("sms", "code", scene, phone)
}

// SmsCodeAttemptKey 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (b *Builder) SmsCodeAttemptKey(scene, phone string) string {
	return b.Key("sms", "code", "attempt", scene, phone)
}

// LoginFailKey 澶勭悊鐧诲綍璁よ瘉骞惰繑鍥炰細璇濈粨鏋溿€
func (b *Builder) LoginFailKey(identifier string) string {
	return b.Key("login", "fail", identifier)
}

// LoginLockKey 澶勭悊鐧诲綍璁よ瘉骞惰繑鍥炰細璇濈粨鏋溿€
func (b *Builder) LoginLockKey(identifier string) string {
	return b.Key("login", "lock", identifier)
}

// MFAChallengeKey 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (b *Builder) MFAChallengeKey(challengeID string) string {
	return b.Key("mfa", "challenge", challengeID)
}

// sanitizePart 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
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

// MustFormat 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func MustFormat(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
