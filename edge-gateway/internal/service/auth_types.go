package service

// AuthRole 表示鉴权后可用的角色条目。
type AuthRole struct {
	// RoleCode 角色业务码（如 SELLER/BUYER），已去除下游枚举前缀。
	RoleCode string
	// ScopeTypeCode 角色作用域类型业务码（如 SHOP/PLATFORM）。
	ScopeTypeCode string
	// ScopeNo 作用域编号字符串，适合前端直接展示/透传。
	ScopeNo string
	// ScopeID 作用域数值主键，供网关内部精细逻辑使用。
	ScopeID uint64
}

// VerifyResult 是网关统一的 token 校验结果。
type VerifyResult struct {
	// UserID 当前访问主体用户 ID。
	UserID uint64
	// AccountStatusCode 账号状态业务码（如 ACTIVE/DISABLED）。
	AccountStatusCode string
	// Roles 当前主体关联的角色集合。
	Roles []AuthRole
}
