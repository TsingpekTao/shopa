package cache

import (
	"context"
	"time"

	"github.com/TsingpekTao/shopa/iam-svc/internal/infra/rediskey"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

// Service 封装 IAM 的 Redis 语义化操作，避免业务层散落 raw key 字符串。
type Service struct {
	keys *rediskey.Builder
}

// New 创建默认 Redis 语义服务。
func New() *Service {
	return &Service{keys: rediskey.New()}
}

// NewWithKeys 允许注入自定义 key builder（测试场景常用）。
func NewWithKeys(keys *rediskey.Builder) *Service {
	return &Service{keys: keys}
}

// IncRegisterIpLimit 对注册 IP 计数并维护窗口 TTL。
func (s *Service) IncRegisterIpLimit(ctx context.Context, ip string, now time.Time) (count int64, err error) {
	key := s.keys.RegIPLimitKey(ip, now)
	count, err = g.Redis().Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	if count == 1 {
		_, _ = g.Redis().Expire(ctx, key, int64(rediskey.TTLRegisterIPWindow/time.Second))
	}
	return count, nil
}

// AcquireSmsSendLock 获取验证码发送锁（60 秒内幂等）。
func (s *Service) AcquireSmsSendLock(ctx context.Context, scene, phone string) (ok bool, err error) {
	key := s.keys.SmsLockKey(scene, phone)
	resp, err := g.Redis().Set(ctx, key, "1", gredis.SetOption{
		NX: true,
		TTLOption: gredis.TTLOption{
			EX: ptrInt64(int64(rediskey.TTLSmsSendLock / time.Second)),
		},
	})
	if err != nil {
		return false, err
	}
	if resp == nil || resp.IsNil() {
		return false, nil
	}
	return true, nil
}

// SaveSmsCode 存储短信码哈希。
func (s *Service) SaveSmsCode(ctx context.Context, scene, phone string, codeHash string) error {
	return g.Redis().SetEX(ctx, s.keys.SmsCodeKey(scene, phone), codeHash, int64(rediskey.TTLSmsCode/time.Second))
}

// GetSmsCodeHash 读取短信码哈希。
func (s *Service) GetSmsCodeHash(ctx context.Context, scene, phone string) (hash string, exists bool, err error) {
	v, err := g.Redis().Get(ctx, s.keys.SmsCodeKey(scene, phone))
	if err != nil {
		return "", false, err
	}
	if v.IsNil() {
		return "", false, nil
	}
	return v.String(), true, nil
}

// DeleteSmsCode 删除短信码及其错误计数。
func (s *Service) DeleteSmsCode(ctx context.Context, scene, phone string) error {
	_, err := g.Redis().Del(ctx, s.keys.SmsCodeKey(scene, phone), s.keys.SmsCodeAttemptKey(scene, phone))
	return err
}

// IncrSmsCodeAttempt 增加短信码错误次数。
func (s *Service) IncrSmsCodeAttempt(ctx context.Context, scene, phone string) (count int64, err error) {
	key := s.keys.SmsCodeAttemptKey(scene, phone)
	count, err = g.Redis().Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	if count == 1 {
		_, _ = g.Redis().Expire(ctx, key, int64(rediskey.TTLSmsCodeAttempt/time.Second))
	}
	return count, nil
}

// IncrLoginFail 增加登录失败计数。
func (s *Service) IncrLoginFail(ctx context.Context, identifier string) (count int64, err error) {
	key := s.keys.LoginFailKey(identifier)
	count, err = g.Redis().Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	if count == 1 {
		_, _ = g.Redis().Expire(ctx, key, int64(rediskey.TTLLoginFailCounter/time.Second))
	}
	return count, nil
}

// IsLoginLocked 查询账号是否在锁定期。
func (s *Service) IsLoginLocked(ctx context.Context, identifier string) (locked bool, err error) {
	n, err := g.Redis().Exists(ctx, s.keys.LoginLockKey(identifier))
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// LockLogin 锁定登录（30 分钟）。
func (s *Service) LockLogin(ctx context.Context, identifier string) error {
	return g.Redis().SetEX(ctx, s.keys.LoginLockKey(identifier), "1", int64(rediskey.TTLLoginLock/time.Second))
}

// ResetLoginFail 重置登录失败计数和锁。
func (s *Service) ResetLoginFail(ctx context.Context, identifier string) error {
	_, err := g.Redis().Del(ctx, s.keys.LoginFailKey(identifier), s.keys.LoginLockKey(identifier))
	return err
}

// SaveMFAChallenge 保存 MFA 挑战。
func (s *Service) SaveMFAChallenge(ctx context.Context, challengeID, payload string) error {
	return g.Redis().SetEX(ctx, s.keys.MFAChallengeKey(challengeID), payload, int64(rediskey.TTLMFAChallenge/time.Second))
}

// GetMFAChallenge 读取 MFA 挑战。
func (s *Service) GetMFAChallenge(ctx context.Context, challengeID string) (payload string, exists bool, err error) {
	v, err := g.Redis().Get(ctx, s.keys.MFAChallengeKey(challengeID))
	if err != nil {
		return "", false, err
	}
	if v == nil || v.IsNil() {
		return "", false, nil
	}
	return v.String(), true, nil
}

// DeleteMFAChallenge 删除 MFA 挑战。
func (s *Service) DeleteMFAChallenge(ctx context.Context, challengeID string) error {
	_, err := g.Redis().Del(ctx, s.keys.MFAChallengeKey(challengeID))
	return err
}

func ptrInt64(v int64) *int64 { return &v }
