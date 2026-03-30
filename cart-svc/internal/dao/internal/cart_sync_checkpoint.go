// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CartSyncCheckpointDao is the data access object for the table cart_sync_checkpoint.
type CartSyncCheckpointDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  CartSyncCheckpointColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// CartSyncCheckpointColumns defines and stores column names for the table cart_sync_checkpoint.
type CartSyncCheckpointColumns struct {
	Id              string //
	UserId          string //
	LastSyncedAt    string //
	LastSyncVersion string //
	CreatedAt       string //
	UpdatedAt       string //
}

// cartSyncCheckpointColumns holds the columns for the table cart_sync_checkpoint.
var cartSyncCheckpointColumns = CartSyncCheckpointColumns{
	Id:              "id",
	UserId:          "user_id",
	LastSyncedAt:    "last_synced_at",
	LastSyncVersion: "last_sync_version",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewCartSyncCheckpointDao creates and returns a new DAO object for table data access.
func NewCartSyncCheckpointDao(handlers ...gdb.ModelHandler) *CartSyncCheckpointDao {
	return &CartSyncCheckpointDao{
		group:    "default",
		table:    "cart_sync_checkpoint",
		columns:  cartSyncCheckpointColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CartSyncCheckpointDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CartSyncCheckpointDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CartSyncCheckpointDao) Columns() CartSyncCheckpointColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CartSyncCheckpointDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CartSyncCheckpointDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CartSyncCheckpointDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
