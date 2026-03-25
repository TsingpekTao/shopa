// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaScenePolicyDao 是 media_scene_policy 表的数据访问对象。
type MediaScenePolicyDao struct {
	table    string                  // table 表示 DAO 所在的表名。
	group    string                  // group 表示数据库配置分组。
	columns  MediaScenePolicyColumns // columns 包含全部列名，方便复用。
	handlers []gdb.ModelHandler      // handlers 用于自定义模型处理。
}

// MediaScenePolicyColumns 定义并存储 media_scene_policy 表的列名。
type MediaScenePolicyColumns struct {
	Id                  string //
	SceneCode           string // 唯一场景码，如 PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC。
	AclType             string // 1 表示 PUBLIC_READ，2 表示 PRIVATE。
	MaxCount            string // 每个绑定槽允许的最大资产数量。
	MaxSizeBytes        string // 文件大小上限（字节），0 表示不限制。
	AllowMimeJson       string // 允许的 MIME 列表 JSON。
	AllowExtJson        string // 允许的扩展名列表 JSON。
	RetentionDays       string // 解绑后保留天数。
	GcGraceHours        string // 物理删除前的额外宽限小时。
	RiskLevel           string // 1 表示 LOW，2 表示 MEDIUM，3 表示 HIGH。
	RiskAsyncEnabled    string // 是否启用异步风控扫描。
	ProcessAsyncEnabled string // 是否启用异步处理流程。
	Status              string // 1 表示 ENABLED，2 表示 DISABLED。
	Remark              string //
	ExtJson             string //
	CreatedAt           string //
	UpdatedAt           string //
}

// mediaScenePolicyColumns 保存 media_scene_policy 表的列名映射。
var mediaScenePolicyColumns = MediaScenePolicyColumns{
	Id:                  "id",
	SceneCode:           "scene_code",
	AclType:             "acl_type",
	MaxCount:            "max_count",
	MaxSizeBytes:        "max_size_bytes",
	AllowMimeJson:       "allow_mime_json",
	AllowExtJson:        "allow_ext_json",
	RetentionDays:       "retention_days",
	GcGraceHours:        "gc_grace_hours",
	RiskLevel:           "risk_level",
	RiskAsyncEnabled:    "risk_async_enabled",
	ProcessAsyncEnabled: "process_async_enabled",
	Status:              "status",
	Remark:              "remark",
	ExtJson:             "ext_json",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewMediaScenePolicyDao 创建并返回 media_scene_policy 表的数据访问对象。
func NewMediaScenePolicyDao(handlers ...gdb.ModelHandler) *MediaScenePolicyDao {
	return &MediaScenePolicyDao{
		group:    "default",
		table:    "media_scene_policy",
		columns:  mediaScenePolicyColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 的底层数据库管理对象。
func (dao *MediaScenePolicyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 操作的表名。
func (dao *MediaScenePolicyDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *MediaScenePolicyDao) Columns() MediaScenePolicyColumns {
	return dao.columns
}

// Group 返回 DAO 使用的数据库配置分组名称。
func (dao *MediaScenePolicyDao) Group() string {
	return dao.group
}

// Ctx 创建并返回当前 DAO 的 Model，并自动设置操作上下文。
func (dao *MediaScenePolicyDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction 用参数函数 f 包裹事务逻辑。
// 如果 f 返回非 nil 错误则回滚并返回该错误。
// 如果 f 返回 nil 则提交并返回 nil。
//
// 注意：f 内请勿显式提交或回滚，事务由此方法自动管理。
func (dao *MediaScenePolicyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
