// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RiskDecisionLogDao is the data access object for the table risk_decision_log.
type RiskDecisionLogDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  RiskDecisionLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// RiskDecisionLogColumns defines and stores column names for the table risk_decision_log.
type RiskDecisionLogColumns struct {
	Id               string //
	DecisionNo       string //
	UserId           string //
	BizType          string //
	BizNo            string //
	Decision         string //
	RiskScore        string //
	MatchedRuleCodes string //
	ReasonCodes      string //
	Degraded         string //
	CreatedAt        string //
}

// riskDecisionLogColumns holds the columns for the table risk_decision_log.
var riskDecisionLogColumns = RiskDecisionLogColumns{
	Id:               "id",
	DecisionNo:       "decision_no",
	UserId:           "user_id",
	BizType:          "biz_type",
	BizNo:            "biz_no",
	Decision:         "decision",
	RiskScore:        "risk_score",
	MatchedRuleCodes: "matched_rule_codes",
	ReasonCodes:      "reason_codes",
	Degraded:         "degraded",
	CreatedAt:        "created_at",
}

// NewRiskDecisionLogDao creates and returns a new DAO object for table data access.
func NewRiskDecisionLogDao(handlers ...gdb.ModelHandler) *RiskDecisionLogDao {
	return &RiskDecisionLogDao{
		group:    "default",
		table:    "risk_decision_log",
		columns:  riskDecisionLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RiskDecisionLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RiskDecisionLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RiskDecisionLogDao) Columns() RiskDecisionLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RiskDecisionLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RiskDecisionLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RiskDecisionLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
