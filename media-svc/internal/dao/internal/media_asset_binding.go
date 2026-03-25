// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaAssetBindingDao is the data access object for the table media_asset_binding.
type MediaAssetBindingDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  MediaAssetBindingColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// MediaAssetBindingColumns defines and stores column names for the table media_asset_binding.
type MediaAssetBindingColumns struct {
	Id              string //
	SceneCode       string //
	BizType         string //
	BizNo           string //
	BindingField    string // Slot name, e.g. main_images/detail_images/carousel_images
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

// mediaAssetBindingColumns holds the columns for the table media_asset_binding.
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

// NewMediaAssetBindingDao creates and returns a new DAO object for table data access.
func NewMediaAssetBindingDao(handlers ...gdb.ModelHandler) *MediaAssetBindingDao {
	return &MediaAssetBindingDao{
		group:    "default",
		table:    "media_asset_binding",
		columns:  mediaAssetBindingColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MediaAssetBindingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MediaAssetBindingDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MediaAssetBindingDao) Columns() MediaAssetBindingColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MediaAssetBindingDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MediaAssetBindingDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MediaAssetBindingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
