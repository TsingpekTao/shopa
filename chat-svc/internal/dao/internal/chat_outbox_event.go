// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChatOutboxEventDao is the data access object for the table chat_outbox_event.
type ChatOutboxEventDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  ChatOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// ChatOutboxEventColumns defines and stores column names for the table chat_outbox_event.
type ChatOutboxEventColumns struct {
	Id           string //
	EventId      string //
	EventType    string //
	AggregateNo  string //
	PayloadJson  string //
	Status       string //
	RetryCount   string //
	NextRetryAt  string //
	PublishedAt  string //
	ErrorMessage string //
	CreatedAt    string //
	UpdatedAt    string //
}

// chatOutboxEventColumns holds the columns for the table chat_outbox_event.
var chatOutboxEventColumns = ChatOutboxEventColumns{
	Id:           "id",
	EventId:      "event_id",
	EventType:    "event_type",
	AggregateNo:  "aggregate_no",
	PayloadJson:  "payload_json",
	Status:       "status",
	RetryCount:   "retry_count",
	NextRetryAt:  "next_retry_at",
	PublishedAt:  "published_at",
	ErrorMessage: "error_message",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewChatOutboxEventDao creates and returns a new DAO object for table data access.
func NewChatOutboxEventDao(handlers ...gdb.ModelHandler) *ChatOutboxEventDao {
	return &ChatOutboxEventDao{
		group:    "default",
		table:    "chat_outbox_event",
		columns:  chatOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ChatOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ChatOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ChatOutboxEventDao) Columns() ChatOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ChatOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ChatOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ChatOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
