package sms

import "context"

const (
	// SceneRegister 注册验证码场景。
	SceneRegister = "register"
	// SceneLogin 登录验证码场景。
	SceneLogin = "login"
	// SceneMFA 二次验证场景。
	SceneMFA = "mfa"
	// SceneResetPassword 重置密码场景。
	SceneResetPassword = "reset_password"
)

// SendCodeRequest 描述发送验证码所需的最小上下文。
type SendCodeRequest struct {
	Scene     string
	Phone     string
	Code      string
	RequestID string
	ClientIP  string
}

// SendCodeResult 封装短信通道返回值，供审计日志落库。
type SendCodeResult struct {
	Provider   string
	BizID      string
	RawCode    string
	RawMessage string
}

// Sender 抽象短信发送器能力。
type Sender interface {
	Provider() string
	SendCode(ctx context.Context, req *SendCodeRequest) (*SendCodeResult, error)
}
