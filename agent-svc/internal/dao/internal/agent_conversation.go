// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AgentConversationDao is the data access object for the table agent_conversation.
type AgentConversationDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  AgentConversationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// AgentConversationColumns defines and stores column names for the table agent_conversation.
type AgentConversationColumns struct {
	Id                      string //
	ConversationNo          string //
	UserId                  string //
	BotCode                 string //
	SceneCode               string //
	ShopNo                  string //
	OrderNo                 string //
	SubOrderNo              string //
	AnchorSpuNo             string //
	AnchorSkuNo             string //
	ConversationStatusCode  string //
	LastRunStatusCode       string //
	IsHumanHandover         string //
	PendingMessageCount     string //
	SessionSummary          string //
	SessionSummaryVersion   string //
	LastSummarizedMessageNo string //
	LastMessageAt           string //
	LastQueueNoticeAt       string //
	CreatedAt               string //
	UpdatedAt               string //
	DeletedAt               string //
}

// agentConversationColumns holds the columns for the table agent_conversation.
var agentConversationColumns = AgentConversationColumns{
	Id:                      "id",
	ConversationNo:          "conversation_no",
	UserId:                  "user_id",
	BotCode:                 "bot_code",
	SceneCode:               "scene_code",
	ShopNo:                  "shop_no",
	OrderNo:                 "order_no",
	SubOrderNo:              "sub_order_no",
	AnchorSpuNo:             "anchor_spu_no",
	AnchorSkuNo:             "anchor_sku_no",
	ConversationStatusCode:  "conversation_status_code",
	LastRunStatusCode:       "last_run_status_code",
	IsHumanHandover:         "is_human_handover",
	PendingMessageCount:     "pending_message_count",
	SessionSummary:          "session_summary",
	SessionSummaryVersion:   "session_summary_version",
	LastSummarizedMessageNo: "last_summarized_message_no",
	LastMessageAt:           "last_message_at",
	LastQueueNoticeAt:       "last_queue_notice_at",
	CreatedAt:               "created_at",
	UpdatedAt:               "updated_at",
	DeletedAt:               "deleted_at",
}

// NewAgentConversationDao creates and returns a new DAO object for table data access.
func NewAgentConversationDao(handlers ...gdb.ModelHandler) *AgentConversationDao {
	return &AgentConversationDao{
		group:    "default",
		table:    "agent_conversation",
		columns:  agentConversationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AgentConversationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AgentConversationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AgentConversationDao) Columns() AgentConversationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AgentConversationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AgentConversationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AgentConversationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
