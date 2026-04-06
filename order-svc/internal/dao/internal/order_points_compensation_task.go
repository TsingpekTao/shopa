// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderPointsCompensationTaskDao is the data access object for the table order_points_compensation_task.
type OrderPointsCompensationTaskDao struct {
	table    string                             // table is the underlying table name of the DAO.
	group    string                             // group is the database configuration group name of the current DAO.
	columns  OrderPointsCompensationTaskColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                 // handlers for customized model modification.
}

// OrderPointsCompensationTaskColumns defines and stores column names for the table order_points_compensation_task.
type OrderPointsCompensationTaskColumns struct {
	Id                  string //
	TaskNo              string //
	OrderNo             string //
	UserId              string //
	PointsReservationNo string //
	ActionCode          string //
	TaskStatus          string //
	RetryCount          string //
	NextRetryAt         string //
	LastError           string //
	CreatedAt           string //
	UpdatedAt           string //
}

// orderPointsCompensationTaskColumns holds the columns for the table order_points_compensation_task.
var orderPointsCompensationTaskColumns = OrderPointsCompensationTaskColumns{
	Id:                  "id",
	TaskNo:              "task_no",
	OrderNo:             "order_no",
	UserId:              "user_id",
	PointsReservationNo: "points_reservation_no",
	ActionCode:          "action_code",
	TaskStatus:          "task_status",
	RetryCount:          "retry_count",
	NextRetryAt:         "next_retry_at",
	LastError:           "last_error",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewOrderPointsCompensationTaskDao creates and returns a new DAO object for table data access.
func NewOrderPointsCompensationTaskDao(handlers ...gdb.ModelHandler) *OrderPointsCompensationTaskDao {
	return &OrderPointsCompensationTaskDao{
		group:    "default",
		table:    "order_points_compensation_task",
		columns:  orderPointsCompensationTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrderPointsCompensationTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrderPointsCompensationTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrderPointsCompensationTaskDao) Columns() OrderPointsCompensationTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrderPointsCompensationTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrderPointsCompensationTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OrderPointsCompensationTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
