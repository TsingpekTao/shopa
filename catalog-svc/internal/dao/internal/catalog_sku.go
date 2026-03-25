// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogSkuDao is the data access object for the table catalog_sku.
type CatalogSkuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CatalogSkuColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CatalogSkuColumns defines and stores column names for the table catalog_sku.
type CatalogSkuColumns struct {
	Id               string //
	SkuNo            string //
	SpuNo            string //
	ShopNo           string //
	SkuName          string //
	SkuImageAssetId  string //
	SkuStatus        string //
	StockStatus      string //
	StockVersion     string // From inventory event
	LastStockEventId string //
	SalePrice        string //
	MarketPrice      string //
	SaleSpecsJson    string // Canonical sale attrs json for display
	SaleSpecsHash    string // Canonical hash for unique combination
	ActiveSpecsHash  string //
	SortOrder        string //
	Version          string //
	DeletedAt        string //
	CreatedAt        string //
	UpdatedAt        string //
}

// catalogSkuColumns holds the columns for the table catalog_sku.
var catalogSkuColumns = CatalogSkuColumns{
	Id:               "id",
	SkuNo:            "sku_no",
	SpuNo:            "spu_no",
	ShopNo:           "shop_no",
	SkuName:          "sku_name",
	SkuImageAssetId:  "sku_image_asset_id",
	SkuStatus:        "sku_status",
	StockStatus:      "stock_status",
	StockVersion:     "stock_version",
	LastStockEventId: "last_stock_event_id",
	SalePrice:        "sale_price",
	MarketPrice:      "market_price",
	SaleSpecsJson:    "sale_specs_json",
	SaleSpecsHash:    "sale_specs_hash",
	ActiveSpecsHash:  "active_specs_hash",
	SortOrder:        "sort_order",
	Version:          "version",
	DeletedAt:        "deleted_at",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
}

// NewCatalogSkuDao creates and returns a new DAO object for table data access.
func NewCatalogSkuDao(handlers ...gdb.ModelHandler) *CatalogSkuDao {
	return &CatalogSkuDao{
		group:    "default",
		table:    "catalog_sku",
		columns:  catalogSkuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogSkuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogSkuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogSkuDao) Columns() CatalogSkuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogSkuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogSkuDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogSkuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
