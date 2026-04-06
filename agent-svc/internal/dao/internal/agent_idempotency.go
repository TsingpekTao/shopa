// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AgentIdempotencyDao is the data access object for the table agent_idempotency.
type AgentIdempotencyDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  AgentIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// AgentIdempotencyColumns defines and stores column names for the table agent_idempotency.
type AgentIdempotencyColumns struct {
	Id             string //
	IdempotencyKey string //
	ScopeCode      string //
	ScopeId        string //
	ResponseJson   string //
	CreatedAt      string //
	UpdatedAt      string //
}

// agentIdempotencyColumns holds the columns for the table agent_idempotency.
var agentIdempotencyColumns = AgentIdempotencyColumns{
	Id:             "id",
	IdempotencyKey: "idempotency_key",
	ScopeCode:      "scope_code",
	ScopeId:        "scope_id",
	ResponseJson:   "response_json",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewAgentIdempotencyDao creates and returns a new DAO object for table data access.
func NewAgentIdempotencyDao(handlers ...gdb.ModelHandler) *AgentIdempotencyDao {
	return &AgentIdempotencyDao{
		group:    "default",
		table:    "agent_idempotency",
		columns:  agentIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AgentIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AgentIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AgentIdempotencyDao) Columns() AgentIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AgentIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AgentIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AgentIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
