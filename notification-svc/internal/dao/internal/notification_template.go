// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationTemplateDao is the data access object for the table notification_template.
type NotificationTemplateDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  NotificationTemplateColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// NotificationTemplateColumns defines and stores column names for the table notification_template.
type NotificationTemplateColumns struct {
	Id              string //
	TemplateCode    string //
	Channel         string //
	BizType         string //
	TitleTemplate   string //
	ContentTemplate string //
	Status          string //
	CreatedAt       string //
	UpdatedAt       string //
}

// notificationTemplateColumns holds the columns for the table notification_template.
var notificationTemplateColumns = NotificationTemplateColumns{
	Id:              "id",
	TemplateCode:    "template_code",
	Channel:         "channel",
	BizType:         "biz_type",
	TitleTemplate:   "title_template",
	ContentTemplate: "content_template",
	Status:          "status",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewNotificationTemplateDao creates and returns a new DAO object for table data access.
func NewNotificationTemplateDao(handlers ...gdb.ModelHandler) *NotificationTemplateDao {
	return &NotificationTemplateDao{
		group:    "default",
		table:    "notification_template",
		columns:  notificationTemplateColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NotificationTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NotificationTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NotificationTemplateDao) Columns() NotificationTemplateColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NotificationTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NotificationTemplateDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *NotificationTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
