// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentEscalationTicket is the golang structure of table agent_escalation_ticket for DAO operations like Where/Data.
type AgentEscalationTicket struct {
	g.Meta                   `orm:"table:agent_escalation_ticket, do:true"`
	Id                       any         //
	TicketNo                 any         //
	ConversationNo           any         //
	UserId                   any         //
	ShopNo                   any         //
	RunNo                    any         //
	EscalationReasonCode     any         //
	StatusCode               any         //
	HandoffSummary           any         //
	HandoffSummaryStatusCode any         //
	OperatorUserId           any         //
	Remark                   any         //
	QueueAppendixJson        any         //
	HandoffGeneratedAt       *gtime.Time //
	AcceptedAt               *gtime.Time //
	ClosedAt                 *gtime.Time //
	CreatedAt                *gtime.Time //
	UpdatedAt                *gtime.Time //
	DeletedAt                *gtime.Time //
}
