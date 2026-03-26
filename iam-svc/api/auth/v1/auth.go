package v1

import "github.com/gogf/gf/v2/frame/g"

// RiskContext 描述请求中的风控上下文信息。
type RiskContext struct {
	Fingerprint  string `json:"fingerprint" dc:"设备指纹"`
	DeviceID     string `json:"deviceId" dc:"设备 ID"`
	CaptchaToken string `json:"captchaToken" dc:"验证码 Token"`
}

// TokenPair 表示 access/refresh 令牌对。
type TokenPair struct {
	TokenType        string `json:"tokenType" dc:"token 类型（当前为 Bearer）"`
	AccessToken      string `json:"accessToken" dc:"access token"`
	AccessExpiresIn  uint32 `json:"accessExpiresIn" dc:"access token 有效期（秒）"`
	RefreshToken     string `json:"refreshToken,omitempty" dc:"refresh token"`
	RefreshExpiresIn uint32 `json:"refreshExpiresIn,omitempty" dc:"refresh token 有效期（秒）"`
	Sid              string `json:"sid" dc:"登录会话 SID"`
}

// MfaChallenge 表示 MFA 二次验证挑战凭据。
type MfaChallenge struct {
	ChallengeID string `json:"challengeId" dc:"MFA 挑战 ID"`
	ExpireAt    string `json:"expireAt,omitempty" dc:"过期时间（RFC3339）"`
}

// AuthResult 使用 oneof 语义返回令牌或 MFA 挑战。
type AuthResult struct {
	TokenPair    *TokenPair    `json:"tokenPair,omitempty"`
	MfaChallenge *MfaChallenge `json:"mfaChallenge,omitempty"`
}

type RoleItem struct {
	RoleCode  int32  `json:"roleCode" dc:"角色编码"`
	ScopeType int32  `json:"scopeType" dc:"作用域类型"`
	ScopeID   uint64 `json:"scopeId" dc:"作用域 ID"`
}

type MembershipSummary struct {
	LevelCode string `json:"levelCode" dc:"会员等级编码"`
	Points    uint64 `json:"points" dc:"会员积分"`
	ExpireAt  string `json:"expireAt,omitempty" dc:"会员过期时间（RFC3339）"`
}

type SessionSummary struct {
	UserID        uint64             `json:"userId"`
	AccountStatus int32              `json:"accountStatus" dc:"账户状态"`
	Roles         []*RoleItem        `json:"roles"`
	Membership    *MembershipSummary `json:"membership"`
	LastLoginAt   string             `json:"lastLoginAt,omitempty" dc:"最近登录时间（RFC3339）"`
	LastLoginIP   string             `json:"lastLoginIp,omitempty"`
}

type SendSmsCodeReq struct {
	g.Meta `path:"/v1/auth/sms/send" method:"post" tags:"Auth" summary:"发送短信验证码，覆盖注册/登录/MFA/重置场景"`
	Scene  int32        `json:"scene" v:"required|in:1,2,3,4" dc:"1=注册 2=登录 3=MFA 4=重置密码"`
	Phone  string       `json:"phone" v:"required|phone" dc:"目标手机号"`
	Risk   *RiskContext `json:"risk" dc:"风控上下文"`
}

type SendSmsCodeRes struct {
	ResendAfterSeconds uint32 `json:"resendAfterSeconds"`
}

type RegisterByPasswordReq struct {
	g.Meta  `path:"/v1/auth/register/password" method:"post" tags:"Auth" summary:"手机号+短信验证码+密码完成注册"`
	Phone   string `json:"phone" v:"required|phone"`
	SmsCode string `json:"smsCode" v:"required|length:4,8"`
	// HTTP 接口使用 `password` 字段；gRPC 内部接口使用 `password_plain` 字段。
	// gRPC RegisterByPassword 接口负责执行密码确认等业务逻辑。
	Password        string       `json:"password" v:"required|length:8,20"`
	ConfirmPassword string       `json:"confirmPassword" v:"required|length:8,20"`
	Risk            *RiskContext `json:"risk"`
}

type RegisterByPasswordRes struct {
	UserID          uint64          `json:"userId"`
	InitDisplayName string          `json:"initDisplayName"`
	Auth            *AuthResult     `json:"auth"`
	Session         *SessionSummary `json:"session"`
}

type LoginByPasswordReq struct {
	g.Meta     `path:"/v1/auth/login/password" method:"post" tags:"Auth" summary:"使用登录标识与密码登录，并在策略触发时要求 MFA"`
	Identifier string       `json:"identifier" v:"required|length:3,128" dc:"登录标识（手机号/邮箱/用户名）"`
	Password   string       `json:"password" v:"required|length:1,64"`
	Risk       *RiskContext `json:"risk"`
}

type LoginByPasswordRes struct {
	Channel int32           `json:"channel"`
	Auth    *AuthResult     `json:"auth"`
	Session *SessionSummary `json:"session"`
}

type LoginBySmsReq struct {
	g.Meta  `path:"/v1/auth/login/sms" method:"post" tags:"Auth" summary:"使用短信验证码登录"`
	Phone   string       `json:"phone" v:"required|phone"`
	SmsCode string       `json:"smsCode" v:"required|length:4,8"`
	Risk    *RiskContext `json:"risk"`
}

type LoginBySmsRes struct {
	Channel int32           `json:"channel"`
	Auth    *AuthResult     `json:"auth"`
	Session *SessionSummary `json:"session"`
}

type VerifyMfaReq struct {
	g.Meta      `path:"/v1/me/mfa/verify" method:"post" tags:"Me" summary:"校验 MFA 挑战并签发最终 token"`
	ChallengeID string       `json:"challengeId" v:"required|length:6,64"`
	SmsCode     string       `json:"smsCode" v:"required|length:4,8"`
	Risk        *RiskContext `json:"risk"`
}

type VerifyMfaRes struct {
	TokenPair *TokenPair      `json:"tokenPair"`
	Session   *SessionSummary `json:"session"`
}

type RefreshTokenReq struct {
	g.Meta       `path:"/v1/auth/token/refresh" method:"post" tags:"Auth" summary:"使用 refresh token 轮换新的 access/refresh 令牌"`
	RefreshToken string       `json:"refreshToken" v:"required" dc:"用于换新的 refresh token"`
	Risk         *RiskContext `json:"risk"`
}

type RefreshTokenRes struct {
	TokenPair *TokenPair `json:"tokenPair"`
}

type LogoutReq struct {
	g.Meta       `path:"/v1/auth/logout" method:"post" tags:"Auth" summary:"注销当前会话或全部会话"`
	RefreshToken string `json:"refreshToken" dc:"待失效的 refresh token（可选）"`
	AllDevices   bool   `json:"allDevices"`
}

type LogoutRes struct{}

type GetMySessionReq struct {
	g.Meta `path:"/v1/me/session" method:"get" tags:"Me" summary:"返回当前 access token 对应的会话摘要"`
}

type GetMySessionRes struct {
	Session *SessionSummary `json:"session"`
}

type ChangePasswordReq struct {
	g.Meta      `path:"/v1/auth/password/change" method:"post" tags:"Auth" summary:"登录态修改密码并使旧令牌失效"`
	OldPassword string       `json:"oldPassword" v:"required|length:1,64"`
	NewPassword string       `json:"newPassword" v:"required|length:8,20"`
	Risk        *RiskContext `json:"risk"`
}

type ChangePasswordRes struct {
	Updated bool `json:"updated"`
}

type ResetPasswordBySmsReq struct {
	g.Meta      `path:"/v1/auth/password/reset/sms" method:"post" tags:"Auth" summary:"通过短信验证码重置密码"`
	Phone       string       `json:"phone" v:"required|phone"`
	SmsCode     string       `json:"smsCode" v:"required|length:4,8"`
	NewPassword string       `json:"newPassword" v:"required|length:8,20"`
	Risk        *RiskContext `json:"risk"`
}

type ResetPasswordBySmsRes struct {
	Updated bool `json:"updated"`
}

type LoginByOAuthReq struct {
	g.Meta   `path:"/v1/auth/login/oauth" method:"post" tags:"Auth" summary:"使用第三方 OAuth code 登录"`
	Provider int32        `json:"provider" v:"required|in:1,2" dc:"OAuth 提供商编号（1/2 等）"`
	Code     string       `json:"code" v:"required"`
	State    string       `json:"state"`
	Risk     *RiskContext `json:"risk"`
}

type LoginByOAuthRes struct {
	Channel int32           `json:"channel"`
	Auth    *AuthResult     `json:"auth"`
	Session *SessionSummary `json:"session"`
}

type BindOAuthReq struct {
	g.Meta   `path:"/v1/auth/oauth/bind" method:"post" tags:"Auth" summary:"绑定第三方 OAuth 账号"`
	Provider int32        `json:"provider" v:"required|in:1,2"`
	Code     string       `json:"code" v:"required"`
	State    string       `json:"state"`
	Risk     *RiskContext `json:"risk"`
}

type BindOAuthRes struct {
	Bound bool `json:"bound"`
}

type UnbindOAuthReq struct {
	g.Meta      `path:"/v1/auth/oauth/unbind" method:"post" tags:"Auth" summary:"解绑第三方 OAuth 账号"`
	Provider    int32        `json:"provider" v:"required|in:1,2"`
	ProviderUID string       `json:"providerUid" v:"required"`
	Risk        *RiskContext `json:"risk"`
}

type UnbindOAuthRes struct {
	Unbound bool `json:"unbound"`
}
