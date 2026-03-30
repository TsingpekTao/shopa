// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RiskEventInboxDao is the data access object for the table risk_event_inbox.
type RiskEventInboxDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  RiskEventInboxColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// RiskEventInboxColumns defines and stores column names for the table risk_event_inbox.
type RiskEventInboxColumns struct {
	Id          string //
	EventId     string //
	EventType   string //
	UserId      string //
	BizNo       string //
	PayloadJson string //
	OccurredAt  string //
	CreatedAt   string //
}

// riskEventInboxColumns holds the columns for the table risk_event_inbox.
var riskEventInboxColumns = RiskEventInboxColumns{
	Id:          "id",
	EventId:     "event_id",
	EventType:   "event_type",
	UserId:      "user_id",
	BizNo:       "biz_no",
	PayloadJson: "payload_json",
	OccurredAt:  "occurred_at",
	CreatedAt:   "created_at",
}

// NewRiskEventInboxDao creates and returns a new DAO object for table data access.
func NewRiskEventInboxDao(handlers ...gdb.ModelHandler) *RiskEventInboxDao {
	return &RiskEventInboxDao{
		group:    "default",
		table:    "risk_event_inbox",
		columns:  riskEventInboxColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RiskEventInboxDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RiskEventInboxDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RiskEventInboxDao) Columns() RiskEventInboxColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RiskEventInboxDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RiskEventInboxDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RiskEventInboxDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
