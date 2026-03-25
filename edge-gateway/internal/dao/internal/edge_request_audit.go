// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EdgeRequestAuditDao is the data access object for the table edge_request_audit.
type EdgeRequestAuditDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  EdgeRequestAuditColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// EdgeRequestAuditColumns defines and stores column names for the table edge_request_audit.
type EdgeRequestAuditColumns struct {
	Id                 string // 涓婚敭ID
	RequestId          string // 璇锋眰ID锛堝叏灞?敮涓?級
	RouteCode          string // 鍛戒腑璺?敱缂栫爜
	UserId             string // 璋冪敤鐢ㄦ埛ID
	Method             string // HTTP鏂规硶
	Path               string // 璇锋眰璺?緞
	UpstreamService    string // 涓婃父鏈嶅姟鏍囪瘑
	StatusCode         string // HTTP鐘舵?鐮
	LatencyMs          string // 璇锋眰鑰楁椂姣??
	Partial            string // 鏄?惁闄嶇骇杩斿洖
	DegradedFieldsJson string // 闄嶇骇瀛楁?鍒楄〃
	ErrorCode          string // 涓氬姟閿欒?鐮
	ClientIp           string // 瀹㈡埛绔疘P
	UserAgent          string // 瀹㈡埛绔疷A
	CreatedAt          string // 鍒涘缓鏃堕棿
}

// edgeRequestAuditColumns holds the columns for the table edge_request_audit.
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
	CreatedAt:          "created_at",
}

// NewEdgeRequestAuditDao creates and returns a new DAO object for table data access.
func NewEdgeRequestAuditDao(handlers ...gdb.ModelHandler) *EdgeRequestAuditDao {
	return &EdgeRequestAuditDao{
		group:    "default",
		table:    "edge_request_audit",
		columns:  edgeRequestAuditColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *EdgeRequestAuditDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *EdgeRequestAuditDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *EdgeRequestAuditDao) Columns() EdgeRequestAuditColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *EdgeRequestAuditDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *EdgeRequestAuditDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *EdgeRequestAuditDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
