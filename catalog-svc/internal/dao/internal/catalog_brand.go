// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogBrandDao is the data access object for the table catalog_brand.
type CatalogBrandDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  CatalogBrandColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// CatalogBrandColumns defines and stores column names for the table catalog_brand.
type CatalogBrandColumns struct {
	Id        string //
	BrandNo   string //
	BrandName string //
	Status    string // 1 ENABLED, 2 DISABLED
	CreatedAt string //
	UpdatedAt string //
}

// catalogBrandColumns holds the columns for the table catalog_brand.
var catalogBrandColumns = CatalogBrandColumns{
	Id:        "id",
	BrandNo:   "brand_no",
	BrandName: "brand_name",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewCatalogBrandDao creates and returns a new DAO object for table data access.
func NewCatalogBrandDao(handlers ...gdb.ModelHandler) *CatalogBrandDao {
	return &CatalogBrandDao{
		group:    "default",
		table:    "catalog_brand",
		columns:  catalogBrandColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogBrandDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogBrandDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogBrandDao) Columns() CatalogBrandColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogBrandDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogBrandDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogBrandDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
