// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionRuleDao is the data access object for the table promotion_rule.
type PromotionRuleDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  PromotionRuleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// PromotionRuleColumns defines and stores column names for the table promotion_rule.
type PromotionRuleColumns struct {
	Id         string //
	RuleCode   string //
	CampaignNo string //
	RuleJson   string //
	Priority   string //
	Status     string //
	CreatedAt  string //
	UpdatedAt  string //
}

// promotionRuleColumns holds the columns for the table promotion_rule.
var promotionRuleColumns = PromotionRuleColumns{
	Id:         "id",
	RuleCode:   "rule_code",
	CampaignNo: "campaign_no",
	RuleJson:   "rule_json",
	Priority:   "priority",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewPromotionRuleDao creates and returns a new DAO object for table data access.
func NewPromotionRuleDao(handlers ...gdb.ModelHandler) *PromotionRuleDao {
	return &PromotionRuleDao{
		group:    "default",
		table:    "promotion_rule",
		columns:  promotionRuleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PromotionRuleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PromotionRuleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PromotionRuleDao) Columns() PromotionRuleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PromotionRuleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PromotionRuleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PromotionRuleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
