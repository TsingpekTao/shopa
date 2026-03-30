// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FulfillmentOperateLogDao is the data access object for the table fulfillment_operate_log.
type FulfillmentOperateLogDao struct {
	table    string                       // table is the underlying table name of the DAO.
	group    string                       // group is the database configuration group name of the current DAO.
	columns  FulfillmentOperateLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler           // handlers for customized model modification.
}

// FulfillmentOperateLogColumns defines and stores column names for the table fulfillment_operate_log.
type FulfillmentOperateLogColumns struct {
	Id           string //
	ShipmentNo   string //
	OperatorType string // SYSTEM/SELLER/INTERNAL
	OperatorId   string //
	Action       string //
	FromStatus   string //
	ToStatus     string //
	Remark       string //
	CreatedAt    string //
}

// fulfillmentOperateLogColumns holds the columns for the table fulfillment_operate_log.
var fulfillmentOperateLogColumns = FulfillmentOperateLogColumns{
	Id:           "id",
	ShipmentNo:   "shipment_no",
	OperatorType: "operator_type",
	OperatorId:   "operator_id",
	Action:       "action",
	FromStatus:   "from_status",
	ToStatus:     "to_status",
	Remark:       "remark",
	CreatedAt:    "created_at",
}

// NewFulfillmentOperateLogDao creates and returns a new DAO object for table data access.
func NewFulfillmentOperateLogDao(handlers ...gdb.ModelHandler) *FulfillmentOperateLogDao {
	return &FulfillmentOperateLogDao{
		group:    "default",
		table:    "fulfillment_operate_log",
		columns:  fulfillmentOperateLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FulfillmentOperateLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FulfillmentOperateLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FulfillmentOperateLogDao) Columns() FulfillmentOperateLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FulfillmentOperateLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FulfillmentOperateLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FulfillmentOperateLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
