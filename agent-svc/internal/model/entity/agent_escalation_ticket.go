// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentEscalationTicket is the golang structure for table agent_escalation_ticket.
type AgentEscalationTicket struct {
	Id                       uint64      `json:"id"                       orm:"id"                          ` //
	TicketNo                 string      `json:"ticketNo"                 orm:"ticket_no"                   ` //
	ConversationNo           string      `json:"conversationNo"           orm:"conversation_no"             ` //
	UserId                   uint64      `json:"userId"                   orm:"user_id"                     ` //
	ShopNo                   string      `json:"shopNo"                   orm:"shop_no"                     ` //
	RunNo                    string      `json:"runNo"                    orm:"run_no"                      ` //
	EscalationReasonCode     string      `json:"escalationReasonCode"     orm:"escalation_reason_code"      ` //
	StatusCode               string      `json:"statusCode"               orm:"status_code"                 ` //
	HandoffSummary           string      `json:"handoffSummary"           orm:"handoff_summary"             ` //
	HandoffSummaryStatusCode string      `json:"handoffSummaryStatusCode" orm:"handoff_summary_status_code" ` //
	OperatorUserId           uint64      `json:"operatorUserId"           orm:"operator_user_id"            ` //
	Remark                   string      `json:"remark"                   orm:"remark"                      ` //
	QueueAppendixJson        string      `json:"queueAppendixJson"        orm:"queue_appendix_json"         ` //
	HandoffGeneratedAt       *gtime.Time `json:"handoffGeneratedAt"       orm:"handoff_generated_at"        ` //
	AcceptedAt               *gtime.Time `json:"acceptedAt"               orm:"accepted_at"                 ` //
	ClosedAt                 *gtime.Time `json:"closedAt"                 orm:"closed_at"                   ` //
	CreatedAt                *gtime.Time `json:"createdAt"                orm:"created_at"                  ` //
	UpdatedAt                *gtime.Time `json:"updatedAt"                orm:"updated_at"                  ` //
	DeletedAt                *gtime.Time `json:"deletedAt"                orm:"deleted_at"                  ` //
}
