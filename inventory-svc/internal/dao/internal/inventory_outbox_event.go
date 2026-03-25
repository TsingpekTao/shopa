// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryOutboxEventDao is the data access object for the table inventory_outbox_event.
type InventoryOutboxEventDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  InventoryOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// InventoryOutboxEventColumns defines and stores column names for the table inventory_outbox_event.
type InventoryOutboxEventColumns struct {
	Id            string //
	EventId       string //
	EventType     string //
	AggregateType string // STOCK/RESERVATION
	AggregateKey  string //
	RequestId     string //
	PayloadJson   string //
	Status        string //
	AvailableAt   string //
	SentAt        string //
	FailCount     string //
	LastError     string //
	CreatedAt     string //
	UpdatedAt     string //
}

// inventoryOutboxEventColumns holds the columns for the table inventory_outbox_event.
var inventoryOutboxEventColumns = InventoryOutboxEventColumns{
	Id:            "id",
	EventId:       "event_id",
	EventType:     "event_type",
	AggregateType: "aggregate_type",
	AggregateKey:  "aggregate_key",
	RequestId:     "request_id",
	PayloadJson:   "payload_json",
	Status:        "status",
	AvailableAt:   "available_at",
	SentAt:        "sent_at",
	FailCount:     "fail_count",
	LastError:     "last_error",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewInventoryOutboxEventDao creates and returns a new DAO object for table data access.
func NewInventoryOutboxEventDao(handlers ...gdb.ModelHandler) *InventoryOutboxEventDao {
	return &InventoryOutboxEventDao{
		group:    "default",
		table:    "inventory_outbox_event",
		columns:  inventoryOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *InventoryOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *InventoryOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *InventoryOutboxEventDao) Columns() InventoryOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *InventoryOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *InventoryOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InventoryOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
