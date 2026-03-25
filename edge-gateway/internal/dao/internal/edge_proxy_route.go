// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EdgeProxyRouteDao is the data access object for the table edge_proxy_route.
type EdgeProxyRouteDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  EdgeProxyRouteColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// EdgeProxyRouteColumns defines and stores column names for the table edge_proxy_route.
type EdgeProxyRouteColumns struct {
	Id                   string // 涓婚敭ID
	RouteCode            string // 璺?敱缂栫爜锛堝叏灞?敮涓?級
	Method               string // HTTP鏂规硶
	PathPattern          string // 缃戝叧璺?敱鍖归厤璺?緞
	UpstreamService      string // 鐩?爣涓婃父鏈嶅姟鏍囪瘑
	UpstreamPathTemplate string // 鐩?爣涓婃父璺?緞妯℃澘
	AuthRequired         string // 鏄?惁闇??閴存潈
	InjectUserContext    string // 鏄?惁娉ㄥ叆鐢ㄦ埛涓婁笅鏂囧ご
	BodyMode             string // 璇锋眰浣撳?鐞嗘ā寮
	TimeoutMs            string // 涓婃父璋冪敤瓒呮椂鏃堕棿锛堟?绉掞級
	RateLimitRps         string // 姣忕?闄愭祦闃堝?锛?琛ㄧず涓嶉檺娴
	Status               string // 鐘舵?锛?鍚?敤 0鍋滅敤
	Remark               string // 澶囨敞
	CreatedAt            string // 鍒涘缓鏃堕棿
	UpdatedAt            string // 鏇存柊鏃堕棿
	DeletedAt            string // 杞?垹闄ゆ椂闂
}

// edgeProxyRouteColumns holds the columns for the table edge_proxy_route.
var edgeProxyRouteColumns = EdgeProxyRouteColumns{
	Id:                   "id",
	RouteCode:            "route_code",
	Method:               "method",
	PathPattern:          "path_pattern",
	UpstreamService:      "upstream_service",
	UpstreamPathTemplate: "upstream_path_template",
	AuthRequired:         "auth_required",
	InjectUserContext:    "inject_user_context",
	BodyMode:             "body_mode",
	TimeoutMs:            "timeout_ms",
	RateLimitRps:         "rate_limit_rps",
	Status:               "status",
	Remark:               "remark",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
	DeletedAt:            "deleted_at",
}

// NewEdgeProxyRouteDao creates and returns a new DAO object for table data access.
func NewEdgeProxyRouteDao(handlers ...gdb.ModelHandler) *EdgeProxyRouteDao {
	return &EdgeProxyRouteDao{
		group:    "default",
		table:    "edge_proxy_route",
		columns:  edgeProxyRouteColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *EdgeProxyRouteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *EdgeProxyRouteDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *EdgeProxyRouteDao) Columns() EdgeProxyRouteColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *EdgeProxyRouteDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *EdgeProxyRouteDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *EdgeProxyRouteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
