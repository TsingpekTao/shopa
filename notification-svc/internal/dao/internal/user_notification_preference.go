// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserNotificationPreferenceDao is the data access object for the table user_notification_preference.
type UserNotificationPreferenceDao struct {
	table    string                            // table is the underlying table name of the DAO.
	group    string                            // group is the database configuration group name of the current DAO.
	columns  UserNotificationPreferenceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                // handlers for customized model modification.
}

// UserNotificationPreferenceColumns defines and stores column names for the table user_notification_preference.
type UserNotificationPreferenceColumns struct {
	Id                      string //
	UserId                  string //
	AllowTransactionalSms   string //
	AllowMarketingSms       string //
	AllowTransactionalEmail string //
	AllowMarketingEmail     string //
	AllowPush               string //
	CreatedAt               string //
	UpdatedAt               string //
}

// userNotificationPreferenceColumns holds the columns for the table user_notification_preference.
var userNotificationPreferenceColumns = UserNotificationPreferenceColumns{
	Id:                      "id",
	UserId:                  "user_id",
	AllowTransactionalSms:   "allow_transactional_sms",
	AllowMarketingSms:       "allow_marketing_sms",
	AllowTransactionalEmail: "allow_transactional_email",
	AllowMarketingEmail:     "allow_marketing_email",
	AllowPush:               "allow_push",
	CreatedAt:               "created_at",
	UpdatedAt:               "updated_at",
}

// NewUserNotificationPreferenceDao creates and returns a new DAO object for table data access.
func NewUserNotificationPreferenceDao(handlers ...gdb.ModelHandler) *UserNotificationPreferenceDao {
	return &UserNotificationPreferenceDao{
		group:    "default",
		table:    "user_notification_preference",
		columns:  userNotificationPreferenceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserNotificationPreferenceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserNotificationPreferenceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserNotificationPreferenceDao) Columns() UserNotificationPreferenceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserNotificationPreferenceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserNotificationPreferenceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserNotificationPreferenceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
