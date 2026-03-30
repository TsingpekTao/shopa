// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SearchOutboxEventDao is the data access object for the table search_outbox_event.
type SearchOutboxEventDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  SearchOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// SearchOutboxEventColumns defines and stores column names for the table search_outbox_event.
type SearchOutboxEventColumns struct {
	Id            string //
	EventId       string //
	AggregateType string //
	AggregateId   string //
	EventType     string //
	PayloadJson   string //
	Status        string //
	RetryCount    string //
	NextRetryAt   string //
	CreatedAt     string //
	UpdatedAt     string //
}

// searchOutboxEventColumns holds the columns for the table search_outbox_event.
var searchOutboxEventColumns = SearchOutboxEventColumns{
	Id:            "id",
	EventId:       "event_id",
	AggregateType: "aggregate_type",
	AggregateId:   "aggregate_id",
	EventType:     "event_type",
	PayloadJson:   "payload_json",
	Status:        "status",
	RetryCount:    "retry_count",
	NextRetryAt:   "next_retry_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewSearchOutboxEventDao creates and returns a new DAO object for table data access.
func NewSearchOutboxEventDao(handlers ...gdb.ModelHandler) *SearchOutboxEventDao {
	return &SearchOutboxEventDao{
		group:    "default",
		table:    "search_outbox_event",
		columns:  searchOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SearchOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SearchOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SearchOutboxEventDao) Columns() SearchOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SearchOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SearchOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SearchOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
