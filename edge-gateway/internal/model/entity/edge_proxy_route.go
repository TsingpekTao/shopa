// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EdgeProxyRoute is the golang structure for table edge_proxy_route.
type EdgeProxyRoute struct {
	Id                    uint64      `json:"id" orm:"id"`
	RouteCode             string      `json:"routeCode" orm:"route_code"`
	Method                string      `json:"method" orm:"method"`
	PathPattern           string      `json:"pathPattern" orm:"path_pattern"`
	UpstreamService       string      `json:"upstreamService" orm:"upstream_service"`
	UpstreamPathTemplate  string      `json:"upstreamPathTemplate" orm:"upstream_path_template"`
	AuthRequired          int         `json:"authRequired" orm:"auth_required"`
	InjectUserContext     int         `json:"injectUserContext" orm:"inject_user_context"`
	BodyMode              string      `json:"bodyMode" orm:"body_mode"`
	TimeoutMs             uint        `json:"timeoutMs" orm:"timeout_ms"`
	RateLimitRps          uint        `json:"rateLimitRps" orm:"rate_limit_rps"`
	Status                uint        `json:"status" orm:"status"`
	Remark                string      `json:"remark" orm:"remark"`
	RequiredPermissionKey string      `json:"requiredPermissionKey" orm:"required_permission_key"`
	Action                string      `json:"action" orm:"action"`
	ResourceIdPathKey     string      `json:"resourceIdPathKey" orm:"resource_id_path_key"`
	CreatedAt             *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt             *gtime.Time `json:"updatedAt" orm:"updated_at"`
	DeletedAt             *gtime.Time `json:"deletedAt" orm:"deleted_at"`
}
