// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsRuleConfigDao is the data access object for the table points_rule_config.
type PointsRuleConfigDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  PointsRuleConfigColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// PointsRuleConfigColumns defines and stores column names for the table points_rule_config.
type PointsRuleConfigColumns struct {
	RuleCode            string //
	RuleName            string //
	StatusCode          string //
	MinOrderAmountCent  string //
	MaxDeductionRateBps string // Max deduction rate in basis points
	DeductPointsPerCent string // How many points are required for one cent discount
	GrantPointsPerCent  string // How many points are granted for one cent paid amount
	RefundGraceDays     string //
	EffectiveAt         string //
	ExpireAt            string //
	RuleSnapshotJson    string //
	CreatedAt           string //
	UpdatedAt           string //
}

// pointsRuleConfigColumns holds the columns for the table points_rule_config.
var pointsRuleConfigColumns = PointsRuleConfigColumns{
	RuleCode:            "rule_code",
	RuleName:            "rule_name",
	StatusCode:          "status_code",
	MinOrderAmountCent:  "min_order_amount_cent",
	MaxDeductionRateBps: "max_deduction_rate_bps",
	DeductPointsPerCent: "deduct_points_per_cent",
	GrantPointsPerCent:  "grant_points_per_cent",
	RefundGraceDays:     "refund_grace_days",
	EffectiveAt:         "effective_at",
	ExpireAt:            "expire_at",
	RuleSnapshotJson:    "rule_snapshot_json",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewPointsRuleConfigDao creates and returns a new DAO object for table data access.
func NewPointsRuleConfigDao(handlers ...gdb.ModelHandler) *PointsRuleConfigDao {
	return &PointsRuleConfigDao{
		group:    "default",
		table:    "points_rule_config",
		columns:  pointsRuleConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsRuleConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsRuleConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsRuleConfigDao) Columns() PointsRuleConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsRuleConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsRuleConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointsRuleConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
