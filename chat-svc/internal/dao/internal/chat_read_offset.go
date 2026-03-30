// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChatReadOffsetDao is the data access object for the table chat_read_offset.
type ChatReadOffsetDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  ChatReadOffsetColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// ChatReadOffsetColumns defines and stores column names for the table chat_read_offset.
type ChatReadOffsetColumns struct {
	Id              string //
	ConversationNo  string //
	ReaderType      string //
	ReaderUserId    string //
	ReadToMessageNo string //
	ReadToMessageId string //
	ReadAt          string //
	CreatedAt       string //
	UpdatedAt       string //
}

// chatReadOffsetColumns holds the columns for the table chat_read_offset.
var chatReadOffsetColumns = ChatReadOffsetColumns{
	Id:              "id",
	ConversationNo:  "conversation_no",
	ReaderType:      "reader_type",
	ReaderUserId:    "reader_user_id",
	ReadToMessageNo: "read_to_message_no",
	ReadToMessageId: "read_to_message_id",
	ReadAt:          "read_at",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewChatReadOffsetDao creates and returns a new DAO object for table data access.
func NewChatReadOffsetDao(handlers ...gdb.ModelHandler) *ChatReadOffsetDao {
	return &ChatReadOffsetDao{
		group:    "default",
		table:    "chat_read_offset",
		columns:  chatReadOffsetColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ChatReadOffsetDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ChatReadOffsetDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ChatReadOffsetDao) Columns() ChatReadOffsetColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ChatReadOffsetDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ChatReadOffsetDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ChatReadOffsetDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
