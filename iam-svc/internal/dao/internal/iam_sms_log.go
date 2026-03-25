// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamSmsLogDao is the data access object for the table iam_sms_log.
type IamSmsLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  IamSmsLogColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// IamSmsLogColumns defines and stores column names for the table iam_sms_log.
type IamSmsLogColumns struct {
	Id          string //
	Scene       string // 1 register,2 login,3 mfa,4 reset_password
	Target      string //
	Provider    string //
	BizId       string //
	Ip          string //
	Ua          string //
	Fingerprint string //
	Success     string //
	ErrorCode   string //
	CreatedAt   string //
}

// iamSmsLogColumns holds the columns for the table iam_sms_log.
var iamSmsLogColumns = IamSmsLogColumns{
	Id:          "id",
	Scene:       "scene",
	Target:      "target",
	Provider:    "provider",
	BizId:       "biz_id",
	Ip:          "ip",
	Ua:          "ua",
	Fingerprint: "fingerprint",
	Success:     "success",
	ErrorCode:   "error_code",
	CreatedAt:   "created_at",
}

// NewIamSmsLogDao creates and returns a new DAO object for table data access.
func NewIamSmsLogDao(handlers ...gdb.ModelHandler) *IamSmsLogDao {
	return &IamSmsLogDao{
		group:    "default",
		table:    "iam_sms_log",
		columns:  iamSmsLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamSmsLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamSmsLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamSmsLogDao) Columns() IamSmsLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamSmsLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamSmsLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamSmsLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
