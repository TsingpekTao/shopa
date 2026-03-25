// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogSkuSaleAttrValueDao is the data access object for the table catalog_sku_sale_attr_value.
type CatalogSkuSaleAttrValueDao struct {
	table    string                         // table is the underlying table name of the DAO.
	group    string                         // group is the database configuration group name of the current DAO.
	columns  CatalogSkuSaleAttrValueColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler             // handlers for customized model modification.
}

// CatalogSkuSaleAttrValueColumns defines and stores column names for the table catalog_sku_sale_attr_value.
type CatalogSkuSaleAttrValueColumns struct {
	Id        string //
	SkuNo     string //
	SpuNo     string //
	AttrCode  string //
	AttrName  string //
	AttrValue string //
	SortOrder string //
	CreatedAt string //
	UpdatedAt string //
}

// catalogSkuSaleAttrValueColumns holds the columns for the table catalog_sku_sale_attr_value.
var catalogSkuSaleAttrValueColumns = CatalogSkuSaleAttrValueColumns{
	Id:        "id",
	SkuNo:     "sku_no",
	SpuNo:     "spu_no",
	AttrCode:  "attr_code",
	AttrName:  "attr_name",
	AttrValue: "attr_value",
	SortOrder: "sort_order",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewCatalogSkuSaleAttrValueDao creates and returns a new DAO object for table data access.
func NewCatalogSkuSaleAttrValueDao(handlers ...gdb.ModelHandler) *CatalogSkuSaleAttrValueDao {
	return &CatalogSkuSaleAttrValueDao{
		group:    "default",
		table:    "catalog_sku_sale_attr_value",
		columns:  catalogSkuSaleAttrValueColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogSkuSaleAttrValueDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogSkuSaleAttrValueDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogSkuSaleAttrValueDao) Columns() CatalogSkuSaleAttrValueColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogSkuSaleAttrValueDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogSkuSaleAttrValueDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogSkuSaleAttrValueDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
