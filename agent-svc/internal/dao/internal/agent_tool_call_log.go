// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AgentToolCallLogDao is the data access object for the table agent_tool_call_log.
type AgentToolCallLogDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  AgentToolCallLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// AgentToolCallLogColumns defines and stores column names for the table agent_tool_call_log.
type AgentToolCallLogColumns struct {
	Id                  string //
	ConversationNo      string //
	RunNo               string //
	ToolName            string //
	AdapterCode         string //
	ToolScopeCode       string //
	SubjectUserId       string //
	SubjectShopNo       string //
	RequestPayloadJson  string //
	ResponsePayloadJson string //
	ToolResultStatus    string //
	DegradedReasonCode  string //
	ErrorCode           string //
	ErrorMessage        string //
	DurationMs          string //
	RequestScopeJson    string //
	SanitizedArgsJson   string //
	ResultCode          string //
	ResultSummary       string //
	LatencyMs           string //
	CreatedAt           string //
	UpdatedAt           string //
}

// agentToolCallLogColumns holds the columns for the table agent_tool_call_log.
var agentToolCallLogColumns = AgentToolCallLogColumns{
	Id:                  "id",
	ConversationNo:      "conversation_no",
	RunNo:               "run_no",
	ToolName:            "tool_name",
	AdapterCode:         "adapter_code",
	ToolScopeCode:       "tool_scope_code",
	SubjectUserId:       "subject_user_id",
	SubjectShopNo:       "subject_shop_no",
	RequestPayloadJson:  "request_payload_json",
	ResponsePayloadJson: "response_payload_json",
	ToolResultStatus:    "tool_result_status",
	DegradedReasonCode:  "degraded_reason_code",
	ErrorCode:           "error_code",
	ErrorMessage:        "error_message",
	DurationMs:          "duration_ms",
	RequestScopeJson:    "request_scope_json",
	SanitizedArgsJson:   "sanitized_args_json",
	ResultCode:          "result_code",
	ResultSummary:       "result_summary",
	LatencyMs:           "latency_ms",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewAgentToolCallLogDao creates and returns a new DAO object for table data access.
func NewAgentToolCallLogDao(handlers ...gdb.ModelHandler) *AgentToolCallLogDao {
	return &AgentToolCallLogDao{
		group:    "default",
		table:    "agent_tool_call_log",
		columns:  agentToolCallLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AgentToolCallLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AgentToolCallLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AgentToolCallLogDao) Columns() AgentToolCallLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AgentToolCallLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AgentToolCallLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AgentToolCallLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
