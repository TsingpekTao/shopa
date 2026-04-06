// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentToolCallLog is the golang structure of table agent_tool_call_log for DAO operations like Where/Data.
type AgentToolCallLog struct {
	g.Meta              `orm:"table:agent_tool_call_log, do:true"`
	Id                  any         //
	ConversationNo      any         //
	RunNo               any         //
	ToolName            any         //
	AdapterCode         any         //
	ToolScopeCode       any         //
	SubjectUserId       any         //
	SubjectShopNo       any         //
	RequestPayloadJson  any         //
	ResponsePayloadJson any         //
	ToolResultStatus    any         //
	DegradedReasonCode  any         //
	ErrorCode           any         //
	ErrorMessage        any         //
	DurationMs          any         //
	RequestScopeJson    any         //
	SanitizedArgsJson   any         //
	ResultCode          any         //
	ResultSummary       any         //
	LatencyMs           any         //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
