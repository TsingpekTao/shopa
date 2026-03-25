// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ConsumerEventDedupDao is the data access object for the table consumer_event_dedup.
type ConsumerEventDedupDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  ConsumerEventDedupColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// ConsumerEventDedupColumns defines and stores column names for the table consumer_event_dedup.
type ConsumerEventDedupColumns struct {
	Id           string // Primary key
	ConsumerName string // Consumer unique name
	EventId      string // Event id for idempotency
	EventType    string // Event type
	ProcessedAt  string // Processed timestamp
	CreatedAt    string // Created timestamp
}

// consumerEventDedupColumns holds the columns for the table consumer_event_dedup.
var consumerEventDedupColumns = ConsumerEventDedupColumns{
	Id:           "id",
	ConsumerName: "consumer_name",
	EventId:      "event_id",
	EventType:    "event_type",
	ProcessedAt:  "processed_at",
	CreatedAt:    "created_at",
}

// NewConsumerEventDedupDao creates and returns a new DAO object for table data access.
func NewConsumerEventDedupDao(handlers ...gdb.ModelHandler) *ConsumerEventDedupDao {
	return &ConsumerEventDedupDao{
		group:    "default",
		table:    "consumer_event_dedup",
		columns:  consumerEventDedupColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ConsumerEventDedupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ConsumerEventDedupDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ConsumerEventDedupDao) Columns() ConsumerEventDedupColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ConsumerEventDedupDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ConsumerEventDedupDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ConsumerEventDedupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
