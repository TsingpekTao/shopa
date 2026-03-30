// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionCampaignDao is the data access object for the table promotion_campaign.
type PromotionCampaignDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  PromotionCampaignColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// PromotionCampaignColumns defines and stores column names for the table promotion_campaign.
type PromotionCampaignColumns struct {
	Id         string //
	CampaignNo string //
	Name       string //
	RuleJson   string //
	StartAt    string //
	EndAt      string //
	Status     string //
	CreatedAt  string //
	UpdatedAt  string //
}

// promotionCampaignColumns holds the columns for the table promotion_campaign.
var promotionCampaignColumns = PromotionCampaignColumns{
	Id:         "id",
	CampaignNo: "campaign_no",
	Name:       "name",
	RuleJson:   "rule_json",
	StartAt:    "start_at",
	EndAt:      "end_at",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewPromotionCampaignDao creates and returns a new DAO object for table data access.
func NewPromotionCampaignDao(handlers ...gdb.ModelHandler) *PromotionCampaignDao {
	return &PromotionCampaignDao{
		group:    "default",
		table:    "promotion_campaign",
		columns:  promotionCampaignColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PromotionCampaignDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PromotionCampaignDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PromotionCampaignDao) Columns() PromotionCampaignColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PromotionCampaignDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PromotionCampaignDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PromotionCampaignDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
