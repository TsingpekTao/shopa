// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryReservationItemDao is the data access object for the table inventory_reservation_item.
type InventoryReservationItemDao struct {
	table    string                          // table is the underlying table name of the DAO.
	group    string                          // group is the database configuration group name of the current DAO.
	columns  InventoryReservationItemColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler              // handlers for customized model modification.
}

// InventoryReservationItemColumns defines and stores column names for the table inventory_reservation_item.
type InventoryReservationItemColumns struct {
	Id            string //
	ReservationNo string //
	OrderNo       string //
	SkuNo         string //
	SpuNo         string //
	ShopNo        string //
	Qty           string //
	StockVersion  string // Version after reserve success
	CreatedAt     string //
	UpdatedAt     string //
}

// inventoryReservationItemColumns holds the columns for the table inventory_reservation_item.
var inventoryReservationItemColumns = InventoryReservationItemColumns{
	Id:            "id",
	ReservationNo: "reservation_no",
	OrderNo:       "order_no",
	SkuNo:         "sku_no",
	SpuNo:         "spu_no",
	ShopNo:        "shop_no",
	Qty:           "qty",
	StockVersion:  "stock_version",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewInventoryReservationItemDao creates and returns a new DAO object for table data access.
func NewInventoryReservationItemDao(handlers ...gdb.ModelHandler) *InventoryReservationItemDao {
	return &InventoryReservationItemDao{
		group:    "default",
		table:    "inventory_reservation_item",
		columns:  inventoryReservationItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *InventoryReservationItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *InventoryReservationItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *InventoryReservationItemDao) Columns() InventoryReservationItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *InventoryReservationItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *InventoryReservationItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InventoryReservationItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
