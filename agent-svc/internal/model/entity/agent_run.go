// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentRun is the golang structure for table agent_run.
type AgentRun struct {
	Id                      uint64      `json:"id"                      orm:"id"                        ` //
	RunNo                   string      `json:"runNo"                   orm:"run_no"                    ` //
	ConversationNo          string      `json:"conversationNo"          orm:"conversation_no"           ` //
	UserId                  uint64      `json:"userId"                  orm:"user_id"                   ` //
	ShopNo                  string      `json:"shopNo"                  orm:"shop_no"                   ` //
	SubjectUserId           uint64      `json:"subjectUserId"           orm:"subject_user_id"           ` //
	SubjectShopNo           string      `json:"subjectShopNo"           orm:"subject_shop_no"           ` //
	TurnNo                  uint64      `json:"turnNo"                  orm:"turn_no"                   ` //
	AcceptedStatusCode      string      `json:"acceptedStatusCode"      orm:"accepted_status_code"      ` //
	RunStatusCode           string      `json:"runStatusCode"           orm:"run_status_code"           ` //
	CurrentNodeCode         string      `json:"currentNodeCode"         orm:"current_node_code"         ` //
	GraphStateJson          string      `json:"graphStateJson"          orm:"graph_state_json"          ` //
	CheckpointVersion       uint64      `json:"checkpointVersion"       orm:"checkpoint_version"        ` //
	ToolWaitTimeoutAt       *gtime.Time `json:"toolWaitTimeoutAt"       orm:"tool_wait_timeout_at"      ` //
	ToolResultStatus        string      `json:"toolResultStatus"        orm:"tool_result_status"        ` //
	DegradedReasonCode      string      `json:"degradedReasonCode"      orm:"degraded_reason_code"      ` //
	QueueBlocked            int         `json:"queueBlocked"            orm:"queue_blocked"             ` //
	QueueHintMessage        string      `json:"queueHintMessage"        orm:"queue_hint_message"        ` //
	RiskDecisionCode        string      `json:"riskDecisionCode"        orm:"risk_decision_code"        ` //
	PromptInjectionFlag     int         `json:"promptInjectionFlag"     orm:"prompt_injection_flag"     ` //
	ReplyInterrupted        int         `json:"replyInterrupted"        orm:"reply_interrupted"         ` //
	InterruptReasonCode     string      `json:"interruptReasonCode"     orm:"interrupt_reason_code"     ` //
	MergedMessageCount      uint        `json:"mergedMessageCount"      orm:"merged_message_count"      ` //
	ToolIterationCount      uint        `json:"toolIterationCount"      orm:"tool_iteration_count"      ` //
	ToolCallCount           uint        `json:"toolCallCount"           orm:"tool_call_count"           ` //
	PromptTokenEstimate     uint        `json:"promptTokenEstimate"     orm:"prompt_token_estimate"     ` //
	CompletionTokenEstimate uint        `json:"completionTokenEstimate" orm:"completion_token_estimate" ` //
	AnswerSourcesJson       string      `json:"answerSourcesJson"       orm:"answer_sources_json"       ` //
	ErrorCode               string      `json:"errorCode"               orm:"error_code"                ` //
	ErrorMessage            string      `json:"errorMessage"            orm:"error_message"             ` //
	CreatedAt               *gtime.Time `json:"createdAt"               orm:"created_at"                ` //
	StartedAt               *gtime.Time `json:"startedAt"               orm:"started_at"                ` //
	FinishedAt              *gtime.Time `json:"finishedAt"              orm:"finished_at"               ` //
	UpdatedAt               *gtime.Time `json:"updatedAt"               orm:"updated_at"                ` //
	DeletedAt               *gtime.Time `json:"deletedAt"               orm:"deleted_at"                ` //
}
