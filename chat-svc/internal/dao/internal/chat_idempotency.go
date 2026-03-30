// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChatIdempotencyDao is the data access object for the table chat_idempotency.
type ChatIdempotencyDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  ChatIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// ChatIdempotencyColumns defines and stores column names for the table chat_idempotency.
type ChatIdempotencyColumns struct {
	Id             string //
	UserId         string //
	BizCode        string //
	IdempotencyKey string //
	TargetNo       string //
	ResponseJson   string //
	Status         string //
	ExpiredAt      string //
	CreatedAt      string //
	UpdatedAt      string //
}

// chatIdempotencyColumns holds the columns for the table chat_idempotency.
var chatIdempotencyColumns = ChatIdempotencyColumns{
	Id:             "id",
	UserId:         "user_id",
	BizCode:        "biz_code",
	IdempotencyKey: "idempotency_key",
	TargetNo:       "target_no",
	ResponseJson:   "response_json",
	Status:         "status",
	ExpiredAt:      "expired_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewChatIdempotencyDao creates and returns a new DAO object for table data access.
func NewChatIdempotencyDao(handlers ...gdb.ModelHandler) *ChatIdempotencyDao {
	return &ChatIdempotencyDao{
		group:    "default",
		table:    "chat_idempotency",
		columns:  chatIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ChatIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ChatIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ChatIdempotencyDao) Columns() ChatIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ChatIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ChatIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ChatIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
