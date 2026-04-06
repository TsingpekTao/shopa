// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AgentRunDao is the data access object for the table agent_run.
type AgentRunDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  AgentRunColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// AgentRunColumns defines and stores column names for the table agent_run.
type AgentRunColumns struct {
	Id                      string //
	RunNo                   string //
	ConversationNo          string //
	UserId                  string //
	ShopNo                  string //
	SubjectUserId           string //
	SubjectShopNo           string //
	TurnNo                  string //
	AcceptedStatusCode      string //
	RunStatusCode           string //
	CurrentNodeCode         string //
	GraphStateJson          string //
	CheckpointVersion       string //
	ToolWaitTimeoutAt       string //
	ToolResultStatus        string //
	DegradedReasonCode      string //
	QueueBlocked            string //
	QueueHintMessage        string //
	RiskDecisionCode        string //
	PromptInjectionFlag     string //
	ReplyInterrupted        string //
	InterruptReasonCode     string //
	MergedMessageCount      string //
	ToolIterationCount      string //
	ToolCallCount           string //
	PromptTokenEstimate     string //
	CompletionTokenEstimate string //
	AnswerSourcesJson       string //
	ErrorCode               string //
	ErrorMessage            string //
	CreatedAt               string //
	StartedAt               string //
	FinishedAt              string //
	UpdatedAt               string //
	DeletedAt               string //
}

// agentRunColumns holds the columns for the table agent_run.
var agentRunColumns = AgentRunColumns{
	Id:                      "id",
	RunNo:                   "run_no",
	ConversationNo:          "conversation_no",
	UserId:                  "user_id",
	ShopNo:                  "shop_no",
	SubjectUserId:           "subject_user_id",
	SubjectShopNo:           "subject_shop_no",
	TurnNo:                  "turn_no",
	AcceptedStatusCode:      "accepted_status_code",
	RunStatusCode:           "run_status_code",
	CurrentNodeCode:         "current_node_code",
	GraphStateJson:          "graph_state_json",
	CheckpointVersion:       "checkpoint_version",
	ToolWaitTimeoutAt:       "tool_wait_timeout_at",
	ToolResultStatus:        "tool_result_status",
	DegradedReasonCode:      "degraded_reason_code",
	QueueBlocked:            "queue_blocked",
	QueueHintMessage:        "queue_hint_message",
	RiskDecisionCode:        "risk_decision_code",
	PromptInjectionFlag:     "prompt_injection_flag",
	ReplyInterrupted:        "reply_interrupted",
	InterruptReasonCode:     "interrupt_reason_code",
	MergedMessageCount:      "merged_message_count",
	ToolIterationCount:      "tool_iteration_count",
	ToolCallCount:           "tool_call_count",
	PromptTokenEstimate:     "prompt_token_estimate",
	CompletionTokenEstimate: "completion_token_estimate",
	AnswerSourcesJson:       "answer_sources_json",
	ErrorCode:               "error_code",
	ErrorMessage:            "error_message",
	CreatedAt:               "created_at",
	StartedAt:               "started_at",
	FinishedAt:              "finished_at",
	UpdatedAt:               "updated_at",
	DeletedAt:               "deleted_at",
}

// NewAgentRunDao creates and returns a new DAO object for table data access.
func NewAgentRunDao(handlers ...gdb.ModelHandler) *AgentRunDao {
	return &AgentRunDao{
		group:    "default",
		table:    "agent_run",
		columns:  agentRunColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AgentRunDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AgentRunDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AgentRunDao) Columns() AgentRunColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AgentRunDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AgentRunDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AgentRunDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
