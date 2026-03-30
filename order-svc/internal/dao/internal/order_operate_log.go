// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderOperateLogDao is the data access object for the table order_operate_log.
type OrderOperateLogDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  OrderOperateLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// OrderOperateLogColumns defines and stores column names for the table order_operate_log.
type OrderOperateLogColumns struct {
	Id               string //
	OrderNo          string //
	SubOrderNo       string //
	OperatorUserId   string //
	OperatorTypeCode string //
	ActionCode       string //
	BeforeStatus     string //
	AfterStatus      string //
	DetailJson       string //
	CreatedAt        string //
}

// orderOperateLogColumns holds the columns for the table order_operate_log.
var orderOperateLogColumns = OrderOperateLogColumns{
	Id:               "id",
	OrderNo:          "order_no",
	SubOrderNo:       "sub_order_no",
	OperatorUserId:   "operator_user_id",
	OperatorTypeCode: "operator_type_code",
	ActionCode:       "action_code",
	BeforeStatus:     "before_status",
	AfterStatus:      "after_status",
	DetailJson:       "detail_json",
	CreatedAt:        "created_at",
}

// NewOrderOperateLogDao creates and returns a new DAO object for table data access.
func NewOrderOperateLogDao(handlers ...gdb.ModelHandler) *OrderOperateLogDao {
	return &OrderOperateLogDao{
		group:    "default",
		table:    "order_operate_log",
		columns:  orderOperateLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrderOperateLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrderOperateLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrderOperateLogDao) Columns() OrderOperateLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrderOperateLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrderOperateLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OrderOperateLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
