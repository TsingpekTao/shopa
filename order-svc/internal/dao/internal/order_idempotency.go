// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderIdempotencyDao is the data access object for the table order_idempotency.
type OrderIdempotencyDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  OrderIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// OrderIdempotencyColumns defines and stores column names for the table order_idempotency.
type OrderIdempotencyColumns struct {
	Id             string //
	UserId         string //
	IdempotencyKey string //
	ActionCode     string //
	OrderNo        string //
	Status         string //
	ResponseJson   string //
	ErrorCode      string //
	ExpireAt       string //
	CreatedAt      string //
	UpdatedAt      string //
}

// orderIdempotencyColumns holds the columns for the table order_idempotency.
var orderIdempotencyColumns = OrderIdempotencyColumns{
	Id:             "id",
	UserId:         "user_id",
	IdempotencyKey: "idempotency_key",
	ActionCode:     "action_code",
	OrderNo:        "order_no",
	Status:         "status",
	ResponseJson:   "response_json",
	ErrorCode:      "error_code",
	ExpireAt:       "expire_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewOrderIdempotencyDao creates and returns a new DAO object for table data access.
func NewOrderIdempotencyDao(handlers ...gdb.ModelHandler) *OrderIdempotencyDao {
	return &OrderIdempotencyDao{
		group:    "default",
		table:    "order_idempotency",
		columns:  orderIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrderIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrderIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrderIdempotencyDao) Columns() OrderIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrderIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrderIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OrderIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
