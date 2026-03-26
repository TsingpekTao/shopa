// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaAssetBindingDao 是 media_asset_binding 表的数据访问对象。
type MediaAssetBindingDao struct {
	table    string                   // table 表示 DAO 底层对应的表名。
	group    string                   // group 表示当前 DAO 的数据库配置分组。
	columns  MediaAssetBindingColumns // columns 包含表的全部列名，方便直接引用。
	handlers []gdb.ModelHandler       // handlers 用于自定义数据模型的修改处理。
}

// MediaAssetBindingColumns 定义并存储 media_asset_binding 表的列名。
type MediaAssetBindingColumns struct {
	Id              string //
	SceneCode       string //
	BizType         string //
	BizNo           string //
	BindingField    string // 插槽名，例如 main_images/detail_images/carousel_images。
	AssetId         string //
	SortOrder       string //
	IsActive        string //
	ActiveAssetSlot string //
	ActiveSortSlot  string //
	OperatorUserId  string //
	RequestId       string //
	UnboundAt       string //
	CreatedAt       string //
	UpdatedAt       string //
}

// mediaAssetBindingColumns 保存 media_asset_binding 表的列名映射。
var mediaAssetBindingColumns = MediaAssetBindingColumns{
	Id:              "id",
	SceneCode:       "scene_code",
	BizType:         "biz_type",
	BizNo:           "biz_no",
	BindingField:    "binding_field",
	AssetId:         "asset_id",
	SortOrder:       "sort_order",
	IsActive:        "is_active",
	ActiveAssetSlot: "active_asset_slot",
	ActiveSortSlot:  "active_sort_slot",
	OperatorUserId:  "operator_user_id",
	RequestId:       "request_id",
	UnboundAt:       "unbound_at",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewMediaAssetBindingDao 创建并返回 media_asset_binding 表的数据访问对象。
func NewMediaAssetBindingDao(handlers ...gdb.ModelHandler) *MediaAssetBindingDao {
	return &MediaAssetBindingDao{
		group:    "default",
		table:    "media_asset_binding",
		columns:  mediaAssetBindingColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 的底层数据库管理对象。
func (dao *MediaAssetBindingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 操作的表名。
func (dao *MediaAssetBindingDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *MediaAssetBindingDao) Columns() MediaAssetBindingColumns {
	return dao.columns
}

// Group 返回 DAO 使用的数据库配置分组名称。
func (dao *MediaAssetBindingDao) Group() string {
	return dao.group
}

// Ctx 创建并返回当前 DAO 的 Model，并自动设置操作上下文。
func (dao *MediaAssetBindingDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MediaAssetBindingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
