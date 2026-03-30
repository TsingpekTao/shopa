// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AftersaleIdempotencyDao is the data access object for the table aftersale_idempotency.
type AftersaleIdempotencyDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  AftersaleIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// AftersaleIdempotencyColumns defines and stores column names for the table aftersale_idempotency.
type AftersaleIdempotencyColumns struct {
	Id             string //
	UserId         string //
	IdempotencyKey string //
	ActionCode     string //
	ResourceNo     string //
	Status         string //
	ResponseJson   string //
	ErrorCode      string //
	ExpireAt       string //
	CreatedAt      string //
	UpdatedAt      string //
}

// aftersaleIdempotencyColumns holds the columns for the table aftersale_idempotency.
var aftersaleIdempotencyColumns = AftersaleIdempotencyColumns{
	Id:             "id",
	UserId:         "user_id",
	IdempotencyKey: "idempotency_key",
	ActionCode:     "action_code",
	ResourceNo:     "resource_no",
	Status:         "status",
	ResponseJson:   "response_json",
	ErrorCode:      "error_code",
	ExpireAt:       "expire_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewAftersaleIdempotencyDao creates and returns a new DAO object for table data access.
func NewAftersaleIdempotencyDao(handlers ...gdb.ModelHandler) *AftersaleIdempotencyDao {
	return &AftersaleIdempotencyDao{
		group:    "default",
		table:    "aftersale_idempotency",
		columns:  aftersaleIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AftersaleIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AftersaleIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AftersaleIdempotencyDao) Columns() AftersaleIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AftersaleIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AftersaleIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AftersaleIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
