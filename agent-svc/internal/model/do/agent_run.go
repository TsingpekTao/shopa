// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentRun is the golang structure of table agent_run for DAO operations like Where/Data.
type AgentRun struct {
	g.Meta                  `orm:"table:agent_run, do:true"`
	Id                      any         //
	RunNo                   any         //
	ConversationNo          any         //
	UserId                  any         //
	ShopNo                  any         //
	SubjectUserId           any         //
	SubjectShopNo           any         //
	TurnNo                  any         //
	AcceptedStatusCode      any         //
	RunStatusCode           any         //
	CurrentNodeCode         any         //
	GraphStateJson          any         //
	CheckpointVersion       any         //
	ToolWaitTimeoutAt       *gtime.Time //
	ToolResultStatus        any         //
	DegradedReasonCode      any         //
	QueueBlocked            any         //
	QueueHintMessage        any         //
	RiskDecisionCode        any         //
	PromptInjectionFlag     any         //
	ReplyInterrupted        any         //
	InterruptReasonCode     any         //
	MergedMessageCount      any         //
	ToolIterationCount      any         //
	ToolCallCount           any         //
	PromptTokenEstimate     any         //
	CompletionTokenEstimate any         //
	AnswerSourcesJson       any         //
	ErrorCode               any         //
	ErrorMessage            any         //
	CreatedAt               *gtime.Time //
	StartedAt               *gtime.Time //
	FinishedAt              *gtime.Time //
	UpdatedAt               *gtime.Time //
	DeletedAt               *gtime.Time //
}
