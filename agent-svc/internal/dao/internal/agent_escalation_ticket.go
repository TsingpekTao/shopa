// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AgentEscalationTicketDao is the data access object for the table agent_escalation_ticket.
type AgentEscalationTicketDao struct {
	table    string                       // table is the underlying table name of the DAO.
	group    string                       // group is the database configuration group name of the current DAO.
	columns  AgentEscalationTicketColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler           // handlers for customized model modification.
}

// AgentEscalationTicketColumns defines and stores column names for the table agent_escalation_ticket.
type AgentEscalationTicketColumns struct {
	Id                       string //
	TicketNo                 string //
	ConversationNo           string //
	UserId                   string //
	ShopNo                   string //
	RunNo                    string //
	EscalationReasonCode     string //
	StatusCode               string //
	HandoffSummary           string //
	HandoffSummaryStatusCode string //
	OperatorUserId           string //
	Remark                   string //
	QueueAppendixJson        string //
	HandoffGeneratedAt       string //
	AcceptedAt               string //
	ClosedAt                 string //
	CreatedAt                string //
	UpdatedAt                string //
	DeletedAt                string //
}

// agentEscalationTicketColumns holds the columns for the table agent_escalation_ticket.
var agentEscalationTicketColumns = AgentEscalationTicketColumns{
	Id:                       "id",
	TicketNo:                 "ticket_no",
	ConversationNo:           "conversation_no",
	UserId:                   "user_id",
	ShopNo:                   "shop_no",
	RunNo:                    "run_no",
	EscalationReasonCode:     "escalation_reason_code",
	StatusCode:               "status_code",
	HandoffSummary:           "handoff_summary",
	HandoffSummaryStatusCode: "handoff_summary_status_code",
	OperatorUserId:           "operator_user_id",
	Remark:                   "remark",
	QueueAppendixJson:        "queue_appendix_json",
	HandoffGeneratedAt:       "handoff_generated_at",
	AcceptedAt:               "accepted_at",
	ClosedAt:                 "closed_at",
	CreatedAt:                "created_at",
	UpdatedAt:                "updated_at",
	DeletedAt:                "deleted_at",
}

// NewAgentEscalationTicketDao creates and returns a new DAO object for table data access.
func NewAgentEscalationTicketDao(handlers ...gdb.ModelHandler) *AgentEscalationTicketDao {
	return &AgentEscalationTicketDao{
		group:    "default",
		table:    "agent_escalation_ticket",
		columns:  agentEscalationTicketColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AgentEscalationTicketDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AgentEscalationTicketDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AgentEscalationTicketDao) Columns() AgentEscalationTicketColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AgentEscalationTicketDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AgentEscalationTicketDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AgentEscalationTicketDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
