// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AftersaleOperateLogDao is the data access object for the table aftersale_operate_log.
type AftersaleOperateLogDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  AftersaleOperateLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// AftersaleOperateLogColumns defines and stores column names for the table aftersale_operate_log.
type AftersaleOperateLogColumns struct {
	Id           string //
	AfterSaleNo  string //
	OperatorType string //
	OperatorId   string //
	ActionCode   string //
	FromStatus   string //
	ToStatus     string //
	ReasonCode   string //
	Remark       string //
	CreatedAt    string //
}

// aftersaleOperateLogColumns holds the columns for the table aftersale_operate_log.
var aftersaleOperateLogColumns = AftersaleOperateLogColumns{
	Id:           "id",
	AfterSaleNo:  "after_sale_no",
	OperatorType: "operator_type",
	OperatorId:   "operator_id",
	ActionCode:   "action_code",
	FromStatus:   "from_status",
	ToStatus:     "to_status",
	ReasonCode:   "reason_code",
	Remark:       "remark",
	CreatedAt:    "created_at",
}

// NewAftersaleOperateLogDao creates and returns a new DAO object for table data access.
func NewAftersaleOperateLogDao(handlers ...gdb.ModelHandler) *AftersaleOperateLogDao {
	return &AftersaleOperateLogDao{
		group:    "default",
		table:    "aftersale_operate_log",
		columns:  aftersaleOperateLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AftersaleOperateLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AftersaleOperateLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AftersaleOperateLogDao) Columns() AftersaleOperateLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AftersaleOperateLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AftersaleOperateLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AftersaleOperateLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
