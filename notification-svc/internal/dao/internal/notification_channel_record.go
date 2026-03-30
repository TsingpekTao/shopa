// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationChannelRecordDao is the data access object for the table notification_channel_record.
type NotificationChannelRecordDao struct {
	table    string                           // table is the underlying table name of the DAO.
	group    string                           // group is the database configuration group name of the current DAO.
	columns  NotificationChannelRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler               // handlers for customized model modification.
}

// NotificationChannelRecordColumns defines and stores column names for the table notification_channel_record.
type NotificationChannelRecordColumns struct {
	Id                string //
	NotificationNo    string //
	ProviderCode      string //
	ProviderMessageId string //
	RequestPayload    string //
	ResponsePayload   string //
	Status            string //
	ErrorCode         string //
	ErrorMessage      string //
	CreatedAt         string //
	UpdatedAt         string //
}

// notificationChannelRecordColumns holds the columns for the table notification_channel_record.
var notificationChannelRecordColumns = NotificationChannelRecordColumns{
	Id:                "id",
	NotificationNo:    "notification_no",
	ProviderCode:      "provider_code",
	ProviderMessageId: "provider_message_id",
	RequestPayload:    "request_payload",
	ResponsePayload:   "response_payload",
	Status:            "status",
	ErrorCode:         "error_code",
	ErrorMessage:      "error_message",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewNotificationChannelRecordDao creates and returns a new DAO object for table data access.
func NewNotificationChannelRecordDao(handlers ...gdb.ModelHandler) *NotificationChannelRecordDao {
	return &NotificationChannelRecordDao{
		group:    "default",
		table:    "notification_channel_record",
		columns:  notificationChannelRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NotificationChannelRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NotificationChannelRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NotificationChannelRecordDao) Columns() NotificationChannelRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NotificationChannelRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NotificationChannelRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *NotificationChannelRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
