// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationTaskDao is the data access object for the table notification_task.
type NotificationTaskDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  NotificationTaskColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// NotificationTaskColumns defines and stores column names for the table notification_task.
type NotificationTaskColumns struct {
	Id                  string //
	NotificationNo      string //
	UserId              string //
	TargetAddress       string //
	TemplateCode        string //
	Channel             string //
	BizType             string //
	TemplateParamsJson  string //
	Status              string //
	BlockedByPreference string //
	BlockedByFrequency  string //
	RetryCount          string //
	NextRetryAt         string //
	CreatedAt           string //
	UpdatedAt           string //
}

// notificationTaskColumns holds the columns for the table notification_task.
var notificationTaskColumns = NotificationTaskColumns{
	Id:                  "id",
	NotificationNo:      "notification_no",
	UserId:              "user_id",
	TargetAddress:       "target_address",
	TemplateCode:        "template_code",
	Channel:             "channel",
	BizType:             "biz_type",
	TemplateParamsJson:  "template_params_json",
	Status:              "status",
	BlockedByPreference: "blocked_by_preference",
	BlockedByFrequency:  "blocked_by_frequency",
	RetryCount:          "retry_count",
	NextRetryAt:         "next_retry_at",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewNotificationTaskDao creates and returns a new DAO object for table data access.
func NewNotificationTaskDao(handlers ...gdb.ModelHandler) *NotificationTaskDao {
	return &NotificationTaskDao{
		group:    "default",
		table:    "notification_task",
		columns:  notificationTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NotificationTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NotificationTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NotificationTaskDao) Columns() NotificationTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NotificationTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NotificationTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *NotificationTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
