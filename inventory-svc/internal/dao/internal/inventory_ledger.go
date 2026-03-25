// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryLedgerDao is the data access object for the table inventory_ledger.
type InventoryLedgerDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  InventoryLedgerColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// InventoryLedgerColumns defines and stores column names for the table inventory_ledger.
type InventoryLedgerColumns struct {
	Id                string //
	TxnNo             string //
	BizType           string //
	BizNo             string //
	ActionCode        string //
	ReservationNo     string //
	OrderNo           string //
	SkuNo             string //
	SpuNo             string //
	ShopNo            string //
	DeltaTotalQty     string //
	DeltaLockedQty    string //
	DeltaAvailableQty string //
	AfterTotalQty     string //
	AfterLockedQty    string //
	AfterAvailableQty string //
	StockVersion      string //
	OperatorType      string // SELLER/ORDER/ADMIN/SYSTEM
	OperatorUserId    string //
	RequestId         string //
	Remark            string //
	CreatedAt         string //
}

// inventoryLedgerColumns holds the columns for the table inventory_ledger.
var inventoryLedgerColumns = InventoryLedgerColumns{
	Id:                "id",
	TxnNo:             "txn_no",
	BizType:           "biz_type",
	BizNo:             "biz_no",
	ActionCode:        "action_code",
	ReservationNo:     "reservation_no",
	OrderNo:           "order_no",
	SkuNo:             "sku_no",
	SpuNo:             "spu_no",
	ShopNo:            "shop_no",
	DeltaTotalQty:     "delta_total_qty",
	DeltaLockedQty:    "delta_locked_qty",
	DeltaAvailableQty: "delta_available_qty",
	AfterTotalQty:     "after_total_qty",
	AfterLockedQty:    "after_locked_qty",
	AfterAvailableQty: "after_available_qty",
	StockVersion:      "stock_version",
	OperatorType:      "operator_type",
	OperatorUserId:    "operator_user_id",
	RequestId:         "request_id",
	Remark:            "remark",
	CreatedAt:         "created_at",
}

// NewInventoryLedgerDao creates and returns a new DAO object for table data access.
func NewInventoryLedgerDao(handlers ...gdb.ModelHandler) *InventoryLedgerDao {
	return &InventoryLedgerDao{
		group:    "default",
		table:    "inventory_ledger",
		columns:  inventoryLedgerColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *InventoryLedgerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *InventoryLedgerDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *InventoryLedgerDao) Columns() InventoryLedgerColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *InventoryLedgerDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *InventoryLedgerDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InventoryLedgerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
