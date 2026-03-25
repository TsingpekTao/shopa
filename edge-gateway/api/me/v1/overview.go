package v1

import "github.com/gogf/gf/v2/frame/g"

type GetOverviewReq struct {
	g.Meta `path:"/v1/me/overview" method:"get" tags:"Me" summary:"Get user overview aggregated by gateway"`
}

type RoleItem struct {
	RoleCode      string `json:"roleCode"`
	ScopeTypeCode string `json:"scopeTypeCode"`
	ScopeNo       string `json:"scopeNo"`
}

type GetOverviewRes struct {
	UserID            uint64     `json:"userId"`
	AccountStatusCode string     `json:"accountStatusCode"`
	Roles             []RoleItem `json:"roles"`
	DisplayName       string     `json:"displayName"`
	AvatarURL         string     `json:"avatarUrl"`
	Points            uint64     `json:"points"`
	Partial           bool       `json:"partial"`
	DegradedFields    []string   `json:"degradedFields"`
}
