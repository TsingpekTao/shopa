// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamOutboxEventArchiveDao is the data access object for the table iam_outbox_event_archive.
type IamOutboxEventArchiveDao struct {
	table    string                       // table is the underlying table name of the DAO.
	group    string                       // group is the database configuration group name of the current DAO.
	columns  IamOutboxEventArchiveColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler           // handlers for customized model modification.
}

// IamOutboxEventArchiveColumns defines and stores column names for the table iam_outbox_event_archive.
type IamOutboxEventArchiveColumns struct {
	Id          string //
	EventId     string //
	EventType   string //
	PayloadJson string //
	Status      string // Archive rows are usually SENT/FAILED/DLQ snapshots
	AvailableAt string //
	SentAt      string //
	FailCount   string //
	LastError   string //
	RetryCount  string //
	NextRetryAt string //
	CreatedAt   string //
	UpdatedAt   string //
	ArchivedAt  string //
}

// iamOutboxEventArchiveColumns holds the columns for the table iam_outbox_event_archive.
var iamOutboxEventArchiveColumns = IamOutboxEventArchiveColumns{
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
	ArchivedAt:  "archived_at",
}

// NewIamOutboxEventArchiveDao creates and returns a new DAO object for table data access.
func NewIamOutboxEventArchiveDao(handlers ...gdb.ModelHandler) *IamOutboxEventArchiveDao {
	return &IamOutboxEventArchiveDao{
		group:    "default",
		table:    "iam_outbox_event_archive",
		columns:  iamOutboxEventArchiveColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamOutboxEventArchiveDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamOutboxEventArchiveDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamOutboxEventArchiveDao) Columns() IamOutboxEventArchiveColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamOutboxEventArchiveDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamOutboxEventArchiveDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamOutboxEventArchiveDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
