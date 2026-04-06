// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ShopStoreCategoryDao is the data access object for the table shop_store_category.
type ShopStoreCategoryDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  ShopStoreCategoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// ShopStoreCategoryColumns defines and stores column names for the table shop_store_category.
type ShopStoreCategoryColumns struct {
	Id           string //
	ShopNo       string //
	ParentId     string //
	Name         string //
	Level        string //
	SortOrder    string //
	IsVisible    string //
	IsDeleted    string //
	ProductCount string //
	CreatedAt    string //
	UpdatedAt    string //
}

// shopStoreCategoryColumns holds the columns for the table shop_store_category.
var shopStoreCategoryColumns = ShopStoreCategoryColumns{
	Id:           "id",
	ShopNo:       "shop_no",
	ParentId:     "parent_id",
	Name:         "name",
	Level:        "level",
	SortOrder:    "sort_order",
	IsVisible:    "is_visible",
	IsDeleted:    "is_deleted",
	ProductCount: "product_count",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewShopStoreCategoryDao creates and returns a new DAO object for table data access.
func NewShopStoreCategoryDao(handlers ...gdb.ModelHandler) *ShopStoreCategoryDao {
	return &ShopStoreCategoryDao{
		group:    "default",
		table:    "shop_store_category",
		columns:  shopStoreCategoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ShopStoreCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ShopStoreCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ShopStoreCategoryDao) Columns() ShopStoreCategoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ShopStoreCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ShopStoreCategoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ShopStoreCategoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
