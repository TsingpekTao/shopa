// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EdgeProxyRouteDao is the data access object for table edge_proxy_route.
type EdgeProxyRouteDao struct {
	table    string
	group    string
	columns  EdgeProxyRouteColumns
	handlers []gdb.ModelHandler
}

// EdgeProxyRouteColumns defines and stores column names for table edge_proxy_route.
type EdgeProxyRouteColumns struct {
	Id                    string
	RouteCode             string
	Method                string
	PathPattern           string
	UpstreamService       string
	UpstreamPathTemplate  string
	AuthRequired          string
	InjectUserContext     string
	BodyMode              string
	TimeoutMs             string
	RateLimitRps          string
	Status                string
	Remark                string
	RequiredPermissionKey string
	Action                string
	ResourceIdPathKey     string
	CreatedAt             string
	UpdatedAt             string
	DeletedAt             string
}

var edgeProxyRouteColumns = EdgeProxyRouteColumns{
	Id:                    "id",
	RouteCode:             "route_code",
	Method:                "method",
	PathPattern:           "path_pattern",
	UpstreamService:       "upstream_service",
	UpstreamPathTemplate:  "upstream_path_template",
	AuthRequired:          "auth_required",
	InjectUserContext:     "inject_user_context",
	BodyMode:              "body_mode",
	TimeoutMs:             "timeout_ms",
	RateLimitRps:          "rate_limit_rps",
	Status:                "status",
	Remark:                "remark",
	RequiredPermissionKey: "required_permission_key",
	Action:                "action",
	ResourceIdPathKey:     "resource_id_path_key",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
	DeletedAt:             "deleted_at",
}

func NewEdgeProxyRouteDao(handlers ...gdb.ModelHandler) *EdgeProxyRouteDao {
	return &EdgeProxyRouteDao{
		group:    "default",
		table:    "edge_proxy_route",
		columns:  edgeProxyRouteColumns,
		handlers: handlers,
	}
}

func (dao *EdgeProxyRouteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

func (dao *EdgeProxyRouteDao) Table() string {
	return dao.table
}

func (dao *EdgeProxyRouteDao) Columns() EdgeProxyRouteColumns {
	return dao.columns
}

func (dao *EdgeProxyRouteDao) Group() string {
	return dao.group
}

func (dao *EdgeProxyRouteDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *EdgeProxyRouteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
