package service

type AuthRole struct {
	RoleCode      string
	ScopeTypeCode string
	ScopeNo       string
	ScopeID       uint64
}

type VerifyResult struct {
	UserID            uint64
	AccountStatusCode string
	Roles             []AuthRole
}
