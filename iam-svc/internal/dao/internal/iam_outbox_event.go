// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamOutboxEventDao is the data access object for the table iam_outbox_event.
type IamOutboxEventDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  IamOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// IamOutboxEventColumns defines and stores column names for the table iam_outbox_event.
type IamOutboxEventColumns struct {
	Id          string //
	EventId     string // Global unique event id
	EventType   string //
	PayloadJson string //
	Status      string // 1 NEW,2 PROCESSING,3 SENT,4 FAILED,5 DLQ
	AvailableAt string // Next dispatch schedule time
	SentAt      string // Successful dispatch time
	FailCount   string // Dispatch failure count
	LastError   string // Last dispatch error text
	RetryCount  string //
	NextRetryAt string //
	CreatedAt   string //
	UpdatedAt   string //
}

// iamOutboxEventColumns holds the columns for the table iam_outbox_event.
var iamOutboxEventColumns = IamOutboxEventColumns{
	Id:          "id",
	EventId:     "event_id",
	EventType:   "event_type",
	PayloadJson: "payload_json",
	Status:      "status",
	AvailableAt: "available_at",
	SentAt:      "sent_at",
	FailCount:   "fail_count",
	LastError:   "last_error",
	RetryCount:  "retry_count",
	NextRetryAt: "next_retry_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewIamOutboxEventDao creates and returns a new DAO object for table data access.
func NewIamOutboxEventDao(handlers ...gdb.ModelHandler) *IamOutboxEventDao {
	return &IamOutboxEventDao{
		group:    "default",
		table:    "iam_outbox_event",
		columns:  iamOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamOutboxEventDao) Columns() IamOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
