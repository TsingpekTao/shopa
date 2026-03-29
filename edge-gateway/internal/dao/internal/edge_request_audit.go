// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EdgeRequestAuditDao is the data access object for table edge_request_audit.
type EdgeRequestAuditDao struct {
	table    string
	group    string
	columns  EdgeRequestAuditColumns
	handlers []gdb.ModelHandler
}

// EdgeRequestAuditColumns defines and stores column names for table edge_request_audit.
type EdgeRequestAuditColumns struct {
	Id                 string
	RequestId          string
	RouteCode          string
	UserId             string
	Method             string
	Path               string
	UpstreamService    string
	StatusCode         string
	LatencyMs          string
	Partial            string
	DegradedFieldsJson string
	ErrorCode          string
	ClientIp           string
	UserAgent          string
	PermissionKey      string
	Action             string
	ResourceId         string
	CreatedAt          string
}

var edgeRequestAuditColumns = EdgeRequestAuditColumns{
	Id:                 "id",
	RequestId:          "request_id",
	RouteCode:          "route_code",
	UserId:             "user_id",
	Method:             "method",
	Path:               "path",
	UpstreamService:    "upstream_service",
	StatusCode:         "status_code",
	LatencyMs:          "latency_ms",
	Partial:            "partial",
	DegradedFieldsJson: "degraded_fields_json",
	ErrorCode:          "error_code",
	ClientIp:           "client_ip",
	UserAgent:          "user_agent",
	PermissionKey:      "permission_key",
	Action:             "action",
	ResourceId:         "resource_id",
	CreatedAt:          "created_at",
}

func NewEdgeRequestAuditDao(handlers ...gdb.ModelHandler) *EdgeRequestAuditDao {
	return &EdgeRequestAuditDao{
		group:    "default",
		table:    "edge_request_audit",
		columns:  edgeRequestAuditColumns,
		handlers: handlers,
	}
}

func (dao *EdgeRequestAuditDao) DB() gdb.DB {
	return g.DB(dao.group)
}

func (dao *EdgeRequestAuditDao) Table() string {
	return dao.table
}

func (dao *EdgeRequestAuditDao) Columns() EdgeRequestAuditColumns {
	return dao.columns
}

func (dao *EdgeRequestAuditDao) Group() string {
	return dao.group
}

func (dao *EdgeRequestAuditDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *EdgeRequestAuditDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
