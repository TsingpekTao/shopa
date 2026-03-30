// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ConversationDao is the data access object for the table conversation.
type ConversationDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  ConversationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// ConversationColumns defines and stores column names for the table conversation.
type ConversationColumns struct {
	Id                 string //
	ConversationNo     string //
	BuyerId            string //
	ShopNo             string //
	SceneCode          string //
	OrderNo            string //
	SubOrderNo         string //
	AnchorSpuNo        string //
	AnchorSkuNo        string //
	LastMessageNo      string //
	LastMessagePreview string //
	LastMessageAt      string //
	ConversationStatus string //
	Version            string //
	CreatedAt          string //
	UpdatedAt          string //
	DeletedAt          string //
}

// conversationColumns holds the columns for the table conversation.
var conversationColumns = ConversationColumns{
	Id:                 "id",
	ConversationNo:     "conversation_no",
	BuyerId:            "buyer_id",
	ShopNo:             "shop_no",
	SceneCode:          "scene_code",
	OrderNo:            "order_no",
	SubOrderNo:         "sub_order_no",
	AnchorSpuNo:        "anchor_spu_no",
	AnchorSkuNo:        "anchor_sku_no",
	LastMessageNo:      "last_message_no",
	LastMessagePreview: "last_message_preview",
	LastMessageAt:      "last_message_at",
	ConversationStatus: "conversation_status",
	Version:            "version",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeletedAt:          "deleted_at",
}

// NewConversationDao creates and returns a new DAO object for table data access.
func NewConversationDao(handlers ...gdb.ModelHandler) *ConversationDao {
	return &ConversationDao{
		group:    "default",
		table:    "conversation",
		columns:  conversationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ConversationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ConversationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ConversationDao) Columns() ConversationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ConversationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ConversationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ConversationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
