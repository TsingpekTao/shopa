// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EdgeProxyRoute is the golang structure of table edge_proxy_route for DAO operations.
type EdgeProxyRoute struct {
	g.Meta                `orm:"table:edge_proxy_route, do:true"`
	Id                    any
	RouteCode             any
	Method                any
	PathPattern           any
	UpstreamService       any
	UpstreamPathTemplate  any
	AuthRequired          any
	InjectUserContext     any
	BodyMode              any
	TimeoutMs             any
	RateLimitRps          any
	Status                any
	Remark                any
	RequiredPermissionKey any
	Action                any
	ResourceIdPathKey     any
	CreatedAt             *gtime.Time
	UpdatedAt             *gtime.Time
	DeletedAt             *gtime.Time
}
