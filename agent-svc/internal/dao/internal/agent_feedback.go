// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AgentFeedbackDao is the data access object for the table agent_feedback.
type AgentFeedbackDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  AgentFeedbackColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// AgentFeedbackColumns defines and stores column names for the table agent_feedback.
type AgentFeedbackColumns struct {
	Id             string //
	ConversationNo string //
	RunNo          string //
	UserId         string //
	ShopNo         string //
	FeedbackCode   string //
	Resolved       string //
	Comment        string //
	CreatedAt      string //
	UpdatedAt      string //
	DeletedAt      string //
}

// agentFeedbackColumns holds the columns for the table agent_feedback.
var agentFeedbackColumns = AgentFeedbackColumns{
	Id:             "id",
	ConversationNo: "conversation_no",
	RunNo:          "run_no",
	UserId:         "user_id",
	ShopNo:         "shop_no",
	FeedbackCode:   "feedback_code",
	Resolved:       "resolved",
	Comment:        "comment",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewAgentFeedbackDao creates and returns a new DAO object for table data access.
func NewAgentFeedbackDao(handlers ...gdb.ModelHandler) *AgentFeedbackDao {
	return &AgentFeedbackDao{
		group:    "default",
		table:    "agent_feedback",
		columns:  agentFeedbackColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AgentFeedbackDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AgentFeedbackDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AgentFeedbackDao) Columns() AgentFeedbackColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AgentFeedbackDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AgentFeedbackDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AgentFeedbackDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
