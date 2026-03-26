// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EdgeProxyRouteDao 是 edge_proxy_route 表的数据访问对象。
type EdgeProxyRouteDao struct {
	table    string                // 表示 DAO 操作对应的底层表名。
	group    string                // 表示当前 DAO 使用的数据库配置分组。
	columns  EdgeProxyRouteColumns // columns 保存 Table 所有列名，便于使用。
	handlers []gdb.ModelHandler    // handlers 用于自定义模型的处理器。
}

// EdgeProxyRouteColumns 定义并保存 edge_proxy_route 表的列名。
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

// edgeProxyRouteColumns 定义 edge_proxy_route 表的列名。
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

// NewEdgeProxyRouteDao 创建并返回 edge_proxy_route 表的数据访问对象。
func NewEdgeProxyRouteDao(handlers ...gdb.ModelHandler) *EdgeProxyRouteDao {
	return &EdgeProxyRouteDao{
		group:    "default",
		table:    "edge_proxy_route",
		columns:  edgeProxyRouteColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 的底层数据库管理对象。
func (dao *EdgeProxyRouteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 操作的表名。
func (dao *EdgeProxyRouteDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *EdgeProxyRouteDao) Columns() EdgeProxyRouteColumns {
	return dao.columns
}

// Group 返回当前 DAO 使用的数据库配置分组名。
func (dao *EdgeProxyRouteDao) Group() string {
	return dao.group
}

// Ctx 创建并返回绑定当前上下文的 Model，自动带上当前操作的上下文。
func (dao *EdgeProxyRouteDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *EdgeProxyRouteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
