// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryStockDao is the data access object for the table inventory_stock.
type InventoryStockDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  InventoryStockColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// InventoryStockColumns defines and stores column names for the table inventory_stock.
type InventoryStockColumns struct {
	Id           string //
	SkuNo        string //
	SpuNo        string //
	ShopNo       string //
	TotalQty     string //
	LockedQty    string //
	AvailableQty string //
	StockVersion string // Monotonic version for projection events
	StockStatus  string //
	IsHot        string //
	RowVersion   string // Internal optimistic lock version
	LastEventId  string //
	CreatedAt    string //
	UpdatedAt    string //
}

// inventoryStockColumns holds the columns for the table inventory_stock.
var inventoryStockColumns = InventoryStockColumns{
	Id:           "id",
	SkuNo:        "sku_no",
	SpuNo:        "spu_no",
	ShopNo:       "shop_no",
	TotalQty:     "total_qty",
	LockedQty:    "locked_qty",
	AvailableQty: "available_qty",
	StockVersion: "stock_version",
	StockStatus:  "stock_status",
	IsHot:        "is_hot",
	RowVersion:   "row_version",
	LastEventId:  "last_event_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewInventoryStockDao creates and returns a new DAO object for table data access.
func NewInventoryStockDao(handlers ...gdb.ModelHandler) *InventoryStockDao {
	return &InventoryStockDao{
		group:    "default",
		table:    "inventory_stock",
		columns:  inventoryStockColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *InventoryStockDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *InventoryStockDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *InventoryStockDao) Columns() InventoryStockColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *InventoryStockDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *InventoryStockDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InventoryStockDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
