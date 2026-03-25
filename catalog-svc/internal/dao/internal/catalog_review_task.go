// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogReviewTaskDao is the data access object for the table catalog_review_task.
type CatalogReviewTaskDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  CatalogReviewTaskColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// CatalogReviewTaskColumns defines and stores column names for the table catalog_review_task.
type CatalogReviewTaskColumns struct {
	Id                 string //
	TaskNo             string //
	SpuNo              string //
	ShopNo             string //
	SpuVersionAtSubmit string //
	SubmitNote         string //
	ReviewStatus       string //
	RejectReasonCode   string //
	RejectComment      string //
	ReviewComment      string //
	ReviewerId         string //
	SubmittedAt        string //
	ReviewStartedAt    string //
	ReviewedAt         string //
	CreatedAt          string //
	UpdatedAt          string //
}

// catalogReviewTaskColumns holds the columns for the table catalog_review_task.
var catalogReviewTaskColumns = CatalogReviewTaskColumns{
	Id:                 "id",
	TaskNo:             "task_no",
	SpuNo:              "spu_no",
	ShopNo:             "shop_no",
	SpuVersionAtSubmit: "spu_version_at_submit",
	SubmitNote:         "submit_note",
	ReviewStatus:       "review_status",
	RejectReasonCode:   "reject_reason_code",
	RejectComment:      "reject_comment",
	ReviewComment:      "review_comment",
	ReviewerId:         "reviewer_id",
	SubmittedAt:        "submitted_at",
	ReviewStartedAt:    "review_started_at",
	ReviewedAt:         "reviewed_at",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
}

// NewCatalogReviewTaskDao creates and returns a new DAO object for table data access.
func NewCatalogReviewTaskDao(handlers ...gdb.ModelHandler) *CatalogReviewTaskDao {
	return &CatalogReviewTaskDao{
		group:    "default",
		table:    "catalog_review_task",
		columns:  catalogReviewTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogReviewTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogReviewTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogReviewTaskDao) Columns() CatalogReviewTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogReviewTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogReviewTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogReviewTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
