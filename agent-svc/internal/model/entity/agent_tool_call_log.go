// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentToolCallLog is the golang structure for table agent_tool_call_log.
type AgentToolCallLog struct {
	Id                  uint64      `json:"id"                  orm:"id"                    ` //
	ConversationNo      string      `json:"conversationNo"      orm:"conversation_no"       ` //
	RunNo               string      `json:"runNo"               orm:"run_no"                ` //
	ToolName            string      `json:"toolName"            orm:"tool_name"             ` //
	AdapterCode         string      `json:"adapterCode"         orm:"adapter_code"          ` //
	ToolScopeCode       string      `json:"toolScopeCode"       orm:"tool_scope_code"       ` //
	SubjectUserId       uint64      `json:"subjectUserId"       orm:"subject_user_id"       ` //
	SubjectShopNo       string      `json:"subjectShopNo"       orm:"subject_shop_no"       ` //
	RequestPayloadJson  string      `json:"requestPayloadJson"  orm:"request_payload_json"  ` //
	ResponsePayloadJson string      `json:"responsePayloadJson" orm:"response_payload_json" ` //
	ToolResultStatus    string      `json:"toolResultStatus"    orm:"tool_result_status"    ` //
	DegradedReasonCode  string      `json:"degradedReasonCode"  orm:"degraded_reason_code"  ` //
	ErrorCode           string      `json:"errorCode"           orm:"error_code"            ` //
	ErrorMessage        string      `json:"errorMessage"        orm:"error_message"         ` //
	DurationMs          uint        `json:"durationMs"          orm:"duration_ms"           ` //
	RequestScopeJson    string      `json:"requestScopeJson"    orm:"request_scope_json"    ` //
	SanitizedArgsJson   string      `json:"sanitizedArgsJson"   orm:"sanitized_args_json"   ` //
	ResultCode          string      `json:"resultCode"          orm:"result_code"           ` //
	ResultSummary       string      `json:"resultSummary"       orm:"result_summary"        ` //
	LatencyMs           uint        `json:"latencyMs"           orm:"latency_ms"            ` //
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            ` //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            ` //
}
