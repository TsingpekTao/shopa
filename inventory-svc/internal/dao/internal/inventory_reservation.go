// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryReservationDao is the data access object for the table inventory_reservation.
type InventoryReservationDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  InventoryReservationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// InventoryReservationColumns defines and stores column names for the table inventory_reservation.
type InventoryReservationColumns struct {
	Id                string //
	ReservationNo     string //
	OrderNo           string //
	UserId            string //
	ReservationStatus string //
	ReserveMode       string //
	ExpiredAt         string //
	ConfirmedAt       string //
	CanceledAt        string //
	CancelReasonCode  string //
	IdempotencyKey    string //
	RequestId         string //
	CreatedAt         string //
	UpdatedAt         string //
}

// inventoryReservationColumns holds the columns for the table inventory_reservation.
var inventoryReservationColumns = InventoryReservationColumns{
	Id:                "id",
	ReservationNo:     "reservation_no",
	OrderNo:           "order_no",
	UserId:            "user_id",
	ReservationStatus: "reservation_status",
	ReserveMode:       "reserve_mode",
	ExpiredAt:         "expired_at",
	ConfirmedAt:       "confirmed_at",
	CanceledAt:        "canceled_at",
	CancelReasonCode:  "cancel_reason_code",
	IdempotencyKey:    "idempotency_key",
	RequestId:         "request_id",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewInventoryReservationDao creates and returns a new DAO object for table data access.
func NewInventoryReservationDao(handlers ...gdb.ModelHandler) *InventoryReservationDao {
	return &InventoryReservationDao{
		group:    "default",
		table:    "inventory_reservation",
		columns:  inventoryReservationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *InventoryReservationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *InventoryReservationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *InventoryReservationDao) Columns() InventoryReservationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *InventoryReservationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *InventoryReservationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InventoryReservationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
