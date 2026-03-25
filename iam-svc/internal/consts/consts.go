package consts

import "time"

const (
	// 账号状态（iam_user_auth.account_status）。
	AccountStatusActive   = 1
	AccountStatusLocked   = 2
	AccountStatusDisabled = 3
	AccountStatusReview   = 4
)

const (
	// 角色编码（iam_user_role.role_code）。
	RoleCodeCustomer = 1
	RoleCodeSeller   = 2
	RoleCodeAdmin    = 3
	RoleCodeCS       = 4
)

const (
	// 作用域类型（iam_user_role.scope_type）。
	ScopeTypeGlobal = 1
	ScopeTypeShop   = 2
)

const (
	// 通用启停状态。
	StatusActive   = 1
	StatusDisabled = 2
)

const (
	// Outbox 状态机。
	OutboxStatusNew        = 1
	OutboxStatusProcessing = 2
	OutboxStatusSent       = 3
	OutboxStatusFailed     = 4
	OutboxStatusDLQ        = 5
)

const (
	// mock 短信通道提供方标识。
	SmsProviderMock = "mock"
	// aliyun 短信通道提供方标识。
	SmsProviderAliyun = "aliyun"
	// 默认会员等级。
	MembershipLevelBasic = "BASIC"
)

const (
	// Token 语义。
	TokenTypeBearer  = "Bearer"
	TokenUseAccess   = "access"
	TokenUseRefresh  = "refresh"
	TokenUseMetadata = "metadata"
)

const (
	// 登录失败锁定阈值。
	DefaultLoginFailMax = 5
)

const (
	// MFA 挑战与锁定默认时长。
	DefaultLockDuration = 30 * time.Minute
	DefaultMfaTTL       = 10 * time.Minute
)
