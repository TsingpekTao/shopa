// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CartOutboxEventDao is the data access object for the table cart_outbox_event.
type CartOutboxEventDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  CartOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// CartOutboxEventColumns defines and stores column names for the table cart_outbox_event.
type CartOutboxEventColumns struct {
	Id            string //
	EventId       string //
	EventType     string //
	AggregateType string //
	AggregateId   string //
	PayloadJson   string //
	HeadersJson   string //
	Status        string // 0=NEW,1=PROCESSING,2=SENT,3=FAILED,4=DLQ
	AvailableAt   string //
	RetryCount    string //
	LastError     string //
	CreatedAt     string //
	UpdatedAt     string //
}

// cartOutboxEventColumns holds the columns for the table cart_outbox_event.
var cartOutboxEventColumns = CartOutboxEventColumns{
	Id:            "id",
	EventId:       "event_id",
	EventType:     "event_type",
	AggregateType: "aggregate_type",
	AggregateId:   "aggregate_id",
	PayloadJson:   "payload_json",
	HeadersJson:   "headers_json",
	Status:        "status",
	AvailableAt:   "available_at",
	RetryCount:    "retry_count",
	LastError:     "last_error",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewCartOutboxEventDao creates and returns a new DAO object for table data access.
func NewCartOutboxEventDao(handlers ...gdb.ModelHandler) *CartOutboxEventDao {
	return &CartOutboxEventDao{
		group:    "default",
		table:    "cart_outbox_event",
		columns:  cartOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CartOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CartOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CartOutboxEventDao) Columns() CartOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CartOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CartOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CartOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
