// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EdgeRequestAudit is the golang structure for table edge_request_audit.
type EdgeRequestAudit struct {
	Id                 uint64      `json:"id" orm:"id"`
	RequestId          string      `json:"requestId" orm:"request_id"`
	RouteCode          string      `json:"routeCode" orm:"route_code"`
	UserId             uint64      `json:"userId" orm:"user_id"`
	Method             string      `json:"method" orm:"method"`
	Path               string      `json:"path" orm:"path"`
	UpstreamService    string      `json:"upstreamService" orm:"upstream_service"`
	StatusCode         int         `json:"statusCode" orm:"status_code"`
	LatencyMs          uint        `json:"latencyMs" orm:"latency_ms"`
	Partial            int         `json:"partial" orm:"partial"`
	DegradedFieldsJson string      `json:"degradedFieldsJson" orm:"degraded_fields_json"`
	ErrorCode          string      `json:"errorCode" orm:"error_code"`
	ClientIp           string      `json:"clientIp" orm:"client_ip"`
	UserAgent          string      `json:"userAgent" orm:"user_agent"`
	PermissionKey      string      `json:"permissionKey" orm:"permission_key"`
	Action             string      `json:"action" orm:"action"`
	ResourceId         string      `json:"resourceId" orm:"resource_id"`
	CreatedAt          *gtime.Time `json:"createdAt" orm:"created_at"`
}
