// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FulfillmentIdempotencyDao is the data access object for the table fulfillment_idempotency.
type FulfillmentIdempotencyDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  FulfillmentIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// FulfillmentIdempotencyColumns defines and stores column names for the table fulfillment_idempotency.
type FulfillmentIdempotencyColumns struct {
	Id             string //
	IdempotencyKey string //
	Scope          string //
	RequestHash    string //
	ResourceId     string //
	ResponseJson   string //
	Status         string // 1 SUCCEEDED 2 PROCESSING 3 FAILED
	ExpireAt       string //
	CreatedAt      string //
	UpdatedAt      string //
}

// fulfillmentIdempotencyColumns holds the columns for the table fulfillment_idempotency.
var fulfillmentIdempotencyColumns = FulfillmentIdempotencyColumns{
	Id:             "id",
	IdempotencyKey: "idempotency_key",
	Scope:          "scope",
	RequestHash:    "request_hash",
	ResourceId:     "resource_id",
	ResponseJson:   "response_json",
	Status:         "status",
	ExpireAt:       "expire_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewFulfillmentIdempotencyDao creates and returns a new DAO object for table data access.
func NewFulfillmentIdempotencyDao(handlers ...gdb.ModelHandler) *FulfillmentIdempotencyDao {
	return &FulfillmentIdempotencyDao{
		group:    "default",
		table:    "fulfillment_idempotency",
		columns:  fulfillmentIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FulfillmentIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FulfillmentIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FulfillmentIdempotencyDao) Columns() FulfillmentIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FulfillmentIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FulfillmentIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FulfillmentIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
