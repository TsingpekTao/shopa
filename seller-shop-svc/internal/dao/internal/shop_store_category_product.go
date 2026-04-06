// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ShopStoreCategoryProductDao is the data access object for the table shop_store_category_product.
type ShopStoreCategoryProductDao struct {
	table    string                          // table is the underlying table name of the DAO.
	group    string                          // group is the database configuration group name of the current DAO.
	columns  ShopStoreCategoryProductColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler              // handlers for customized model modification.
}

// ShopStoreCategoryProductColumns defines and stores column names for the table shop_store_category_product.
type ShopStoreCategoryProductColumns struct {
	Id              string //
	ShopNo          string //
	SpuNo           string //
	StoreCategoryId string //
	CreatedAt       string //
	UpdatedAt       string //
}

// shopStoreCategoryProductColumns holds the columns for the table shop_store_category_product.
var shopStoreCategoryProductColumns = ShopStoreCategoryProductColumns{
	Id:              "id",
	ShopNo:          "shop_no",
	SpuNo:           "spu_no",
	StoreCategoryId: "store_category_id",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewShopStoreCategoryProductDao creates and returns a new DAO object for table data access.
func NewShopStoreCategoryProductDao(handlers ...gdb.ModelHandler) *ShopStoreCategoryProductDao {
	return &ShopStoreCategoryProductDao{
		group:    "default",
		table:    "shop_store_category_product",
		columns:  shopStoreCategoryProductColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ShopStoreCategoryProductDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ShopStoreCategoryProductDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ShopStoreCategoryProductDao) Columns() ShopStoreCategoryProductColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ShopStoreCategoryProductDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ShopStoreCategoryProductDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ShopStoreCategoryProductDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
