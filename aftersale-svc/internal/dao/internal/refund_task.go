// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RefundTaskDao is the data access object for the table refund_task.
type RefundTaskDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RefundTaskColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RefundTaskColumns defines and stores column names for the table refund_task.
type RefundTaskColumns struct {
	Id               string //
	RefundTaskNo     string //
	AfterSaleNo      string //
	OrderNo          string //
	SubOrderNo       string //
	PayNo            string //
	RefundAmount     string //
	Status           string //
	RetryCount       string //
	NextRetryAt      string //
	LastErrorCode    string //
	LastErrorMessage string //
	Version          string //
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// refundTaskColumns holds the columns for the table refund_task.
var refundTaskColumns = RefundTaskColumns{
	Id:               "id",
	RefundTaskNo:     "refund_task_no",
	AfterSaleNo:      "after_sale_no",
	OrderNo:          "order_no",
	SubOrderNo:       "sub_order_no",
	PayNo:            "pay_no",
	RefundAmount:     "refund_amount",
	Status:           "status",
	RetryCount:       "retry_count",
	NextRetryAt:      "next_retry_at",
	LastErrorCode:    "last_error_code",
	LastErrorMessage: "last_error_message",
	Version:          "version",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewRefundTaskDao creates and returns a new DAO object for table data access.
func NewRefundTaskDao(handlers ...gdb.ModelHandler) *RefundTaskDao {
	return &RefundTaskDao{
		group:    "default",
		table:    "refund_task",
		columns:  refundTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RefundTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RefundTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RefundTaskDao) Columns() RefundTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RefundTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RefundTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RefundTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
