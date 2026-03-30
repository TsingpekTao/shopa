// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewOutboxEventDao is the data access object for the table review_outbox_event.
type ReviewOutboxEventDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  ReviewOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// ReviewOutboxEventColumns defines and stores column names for the table review_outbox_event.
type ReviewOutboxEventColumns struct {
	Id            string //
	EventId       string //
	AggregateType string //
	AggregateId   string //
	EventType     string //
	PayloadJson   string //
	Status        string //
	RetryCount    string //
	NextRetryAt   string //
	PublishedAt   string //
	CreatedAt     string //
	UpdatedAt     string //
}

// reviewOutboxEventColumns holds the columns for the table review_outbox_event.
var reviewOutboxEventColumns = ReviewOutboxEventColumns{
	Id:            "id",
	EventId:       "event_id",
	AggregateType: "aggregate_type",
	AggregateId:   "aggregate_id",
	EventType:     "event_type",
	PayloadJson:   "payload_json",
	Status:        "status",
	RetryCount:    "retry_count",
	NextRetryAt:   "next_retry_at",
	PublishedAt:   "published_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewReviewOutboxEventDao creates and returns a new DAO object for table data access.
func NewReviewOutboxEventDao(handlers ...gdb.ModelHandler) *ReviewOutboxEventDao {
	return &ReviewOutboxEventDao{
		group:    "default",
		table:    "review_outbox_event",
		columns:  reviewOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReviewOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReviewOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReviewOutboxEventDao) Columns() ReviewOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReviewOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReviewOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReviewOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
