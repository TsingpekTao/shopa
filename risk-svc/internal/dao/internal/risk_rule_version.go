// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RiskRuleVersionDao is the data access object for the table risk_rule_version.
type RiskRuleVersionDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  RiskRuleVersionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// RiskRuleVersionColumns defines and stores column names for the table risk_rule_version.
type RiskRuleVersionColumns struct {
	Id           string //
	RuleCode     string //
	Version      string //
	RuleExprJson string //
	Enabled      string //
	CreatedAt    string //
}

// riskRuleVersionColumns holds the columns for the table risk_rule_version.
var riskRuleVersionColumns = RiskRuleVersionColumns{
	Id:           "id",
	RuleCode:     "rule_code",
	Version:      "version",
	RuleExprJson: "rule_expr_json",
	Enabled:      "enabled",
	CreatedAt:    "created_at",
}

// NewRiskRuleVersionDao creates and returns a new DAO object for table data access.
func NewRiskRuleVersionDao(handlers ...gdb.ModelHandler) *RiskRuleVersionDao {
	return &RiskRuleVersionDao{
		group:    "default",
		table:    "risk_rule_version",
		columns:  riskRuleVersionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RiskRuleVersionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RiskRuleVersionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RiskRuleVersionDao) Columns() RiskRuleVersionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RiskRuleVersionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RiskRuleVersionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RiskRuleVersionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
