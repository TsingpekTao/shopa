// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EdgeRequestAuditDao 是 edge_request_audit 表的数据访问对象。
type EdgeRequestAuditDao struct {
	table    string                  // 表示 DAO 操作对应的底层表名。
	group    string                  // 表示当前 DAO 使用的数据库配置分组。
	columns  EdgeRequestAuditColumns // columns 保存 Table 所有列名，便于使用。
	handlers []gdb.ModelHandler      // handlers 用于自定义模型的处理器。
}

// EdgeRequestAuditColumns 定义并保存 edge_request_audit 表的列名。
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

// edgeRequestAuditColumns 定义 edge_request_audit 表的列名。
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

// NewEdgeRequestAuditDao 创建并返回 edge_request_audit 表的数据访问对象。
func NewEdgeRequestAuditDao(handlers ...gdb.ModelHandler) *EdgeRequestAuditDao {
	return &EdgeRequestAuditDao{
		group:    "default",
		table:    "edge_request_audit",
		columns:  edgeRequestAuditColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 的底层数据库管理对象。
func (dao *EdgeRequestAuditDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 操作的表名。
func (dao *EdgeRequestAuditDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *EdgeRequestAuditDao) Columns() EdgeRequestAuditColumns {
	return dao.columns
}

// Group 返回当前 DAO 使用的数据库配置分组名。
func (dao *EdgeRequestAuditDao) Group() string {
	return dao.group
}

// Ctx 创建并返回绑定当前上下文的 Model，自动带上当前操作的上下文。
func (dao *EdgeRequestAuditDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction 使用函数 f 封装事务逻辑。
// 如果 f 返回非 nil 错误，则回滚事务并返回该错误。
// 如果 f 返回 nil，则提交事务并返回 nil。
//
// 注意：不要在 f 内显式提交或回滚事务，
// 因为本函数会自动处理。
func (dao *EdgeRequestAuditDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
