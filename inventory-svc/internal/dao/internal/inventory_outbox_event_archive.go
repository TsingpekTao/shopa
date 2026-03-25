// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryOutboxEventArchiveDao is the data access object for the table inventory_outbox_event_archive.
type InventoryOutboxEventArchiveDao struct {
	table    string                             // table is the underlying table name of the DAO.
	group    string                             // group is the database configuration group name of the current DAO.
	columns  InventoryOutboxEventArchiveColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                 // handlers for customized model modification.
}

// InventoryOutboxEventArchiveColumns defines and stores column names for the table inventory_outbox_event_archive.
type InventoryOutboxEventArchiveColumns struct {
	Id            string //
	EventId       string //
	EventType     string //
	AggregateType string //
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
	ArchivedAt    string //
}

// inventoryOutboxEventArchiveColumns holds the columns for the table inventory_outbox_event_archive.
var inventoryOutboxEventArchiveColumns = InventoryOutboxEventArchiveColumns{
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
	ArchivedAt:    "archived_at",
}

// NewInventoryOutboxEventArchiveDao creates and returns a new DAO object for table data access.
func NewInventoryOutboxEventArchiveDao(handlers ...gdb.ModelHandler) *InventoryOutboxEventArchiveDao {
	return &InventoryOutboxEventArchiveDao{
		group:    "default",
		table:    "inventory_outbox_event_archive",
		columns:  inventoryOutboxEventArchiveColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *InventoryOutboxEventArchiveDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *InventoryOutboxEventArchiveDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *InventoryOutboxEventArchiveDao) Columns() InventoryOutboxEventArchiveColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *InventoryOutboxEventArchiveDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *InventoryOutboxEventArchiveDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InventoryOutboxEventArchiveDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
