// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EdgeRequestAudit is the golang structure of table edge_request_audit for DAO operations.
type EdgeRequestAudit struct {
	g.Meta             `orm:"table:edge_request_audit, do:true"`
	Id                 any
	RequestId          any
	RouteCode          any
	UserId             any
	Method             any
	Path               any
	UpstreamService    any
	StatusCode         any
	LatencyMs          any
	Partial            any
	DegradedFieldsJson any
	ErrorCode          any
	ClientIp           any
	UserAgent          any
	PermissionKey      any
	Action             any
	ResourceId         any
	CreatedAt          *gtime.Time
}
