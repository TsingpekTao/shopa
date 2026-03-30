// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationDeadLetterDao is the data access object for the table notification_dead_letter.
type NotificationDeadLetterDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  NotificationDeadLetterColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// NotificationDeadLetterColumns defines and stores column names for the table notification_dead_letter.
type NotificationDeadLetterColumns struct {
	Id             string //
	NotificationNo string //
	ReasonCode     string //
	PayloadJson    string //
	CreatedAt      string //
}

// notificationDeadLetterColumns holds the columns for the table notification_dead_letter.
var notificationDeadLetterColumns = NotificationDeadLetterColumns{
	Id:             "id",
	NotificationNo: "notification_no",
	ReasonCode:     "reason_code",
	PayloadJson:    "payload_json",
	CreatedAt:      "created_at",
}

// NewNotificationDeadLetterDao creates and returns a new DAO object for table data access.
func NewNotificationDeadLetterDao(handlers ...gdb.ModelHandler) *NotificationDeadLetterDao {
	return &NotificationDeadLetterDao{
		group:    "default",
		table:    "notification_dead_letter",
		columns:  notificationDeadLetterColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NotificationDeadLetterDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NotificationDeadLetterDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NotificationDeadLetterDao) Columns() NotificationDeadLetterColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NotificationDeadLetterDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NotificationDeadLetterDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *NotificationDeadLetterDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
