package rediskey

import (
	"testing"
	"time"
)

// TestBuilderPrefix 验证 Builder 返回的前缀包含 env/project/service。
func TestBuilderPrefix(t *testing.T) {
	b := NewWith("dev", "shopa", "iam")
	if got, want := b.Prefix(), "dev:shopa:iam"; got != want {
		t.Fatalf("unexpected prefix: got=%s want=%s", got, want)
	}
}

// TestBuilderNoConflictAcrossEnvAndService 验证不同 env/service 不会冲突。
func TestBuilderNoConflictAcrossEnvAndService(t *testing.T) {
	tm := time.Date(2026, 3, 22, 21, 0, 0, 0, time.UTC)
	k1 := NewWith("dev", "shopa", "iam").RegIPLimitKey("1.2.3.4", tm)
	k2 := NewWith("prod", "shopa", "iam").RegIPLimitKey("1.2.3.4", tm)
	k3 := NewWith("dev", "shopa", "order").RegIPLimitKey("1.2.3.4", tm)

	if k1 == k2 {
		t.Fatalf("env collision: %s == %s", k1, k2)
	}
	if k1 == k3 {
		t.Fatalf("service collision: %s == %s", k1, k3)
	}
}

// TestBuilderExpectedKeys 验证各类 Redis key 都按照预期格式生成。
func TestBuilderExpectedKeys(t *testing.T) {
	b := NewWith("dev", "shopa", "iam")
	tm := time.Date(2026, 3, 22, 21, 0, 0, 0, time.UTC)

	cases := map[string]string{
		"reg_ip":        b.RegIPLimitKey("10.0.0.1", tm),
		"sms_lock":      b.SmsLockKey("REGISTER", "13800000000"),
		"sms_code":      b.SmsCodeKey("LOGIN", "13800000000"),
		"login_fail":    b.LoginFailKey("phone_13800000000"),
		"login_lock":    b.LoginLockKey("phone_13800000000"),
		"mfa_challenge": b.MFAChallengeKey("ch_123"),
	}

	expected := map[string]string{
		"reg_ip":        "dev:shopa:iam:reg:ip:10.0.0.1:2026032221",
		"sms_lock":      "dev:shopa:iam:sms:lock:REGISTER:13800000000",
		"sms_code":      "dev:shopa:iam:sms:code:LOGIN:13800000000",
		"login_fail":    "dev:shopa:iam:login:fail:phone_13800000000",
		"login_lock":    "dev:shopa:iam:login:lock:phone_13800000000",
		"mfa_challenge": "dev:shopa:iam:mfa:challenge:ch_123",
	}

	for name, got := range cases {
		if want := expected[name]; got != want {
			t.Fatalf("%s key mismatch: got=%s want=%s", name, got, want)
		}
	}
}
