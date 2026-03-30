// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewIdempotencyDao is the data access object for the table review_idempotency.
type ReviewIdempotencyDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  ReviewIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// ReviewIdempotencyColumns defines and stores column names for the table review_idempotency.
type ReviewIdempotencyColumns struct {
	Id             string //
	UserId         string //
	Action         string //
	IdempotencyKey string //
	ResourceNo     string //
	Status         string //
	ResponseJson   string //
	ErrorCode      string //
	CreatedAt      string //
	UpdatedAt      string //
}

// reviewIdempotencyColumns holds the columns for the table review_idempotency.
var reviewIdempotencyColumns = ReviewIdempotencyColumns{
	Id:             "id",
	UserId:         "user_id",
	Action:         "action",
	IdempotencyKey: "idempotency_key",
	ResourceNo:     "resource_no",
	Status:         "status",
	ResponseJson:   "response_json",
	ErrorCode:      "error_code",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewReviewIdempotencyDao creates and returns a new DAO object for table data access.
func NewReviewIdempotencyDao(handlers ...gdb.ModelHandler) *ReviewIdempotencyDao {
	return &ReviewIdempotencyDao{
		group:    "default",
		table:    "review_idempotency",
		columns:  reviewIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReviewIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReviewIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReviewIdempotencyDao) Columns() ReviewIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReviewIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReviewIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReviewIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
