// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SellerApplicationDao is the data access object for the table seller_application.
type SellerApplicationDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  SellerApplicationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// SellerApplicationColumns defines and stores column names for the table seller_application.
type SellerApplicationColumns struct {
	Id                    string //
	ApplicationNo         string // External application id
	Version               string // Optimistic version
	PreviousApplicationNo string // Chain link to previous rejected application
	NextApplicationNo     string // Chain link to next resubmitted application
	OwnerUserId           string //
	Status                string // ApplicationStatus enum
	EntityNo              string //
	ShopNo                string //
	EntityName            string // Denormalized for list/search
	ShopName              string // Denormalized for list/search
	ShopNameNorm          string // Denormalized normalized shop name
	EntityDraftJson       string // Current draft entity payload
	ShopDraftJson         string // Current draft shop payload
	EntitySubmittedJson   string // Snapshot for review at submit time
	ShopSubmittedJson     string // Snapshot for review at submit time
	LatestRejectJson      string // Latest reject info
	SubmittedAt           string //
	ReviewStartedAt       string //
	ReviewedAt            string //
	ReviewerId            string //
	ReviewComment         string //
	CreatedAt             string //
	UpdatedAt             string //
}

// sellerApplicationColumns holds the columns for the table seller_application.
var sellerApplicationColumns = SellerApplicationColumns{
	Id:                    "id",
	ApplicationNo:         "application_no",
	Version:               "version",
	PreviousApplicationNo: "previous_application_no",
	NextApplicationNo:     "next_application_no",
	OwnerUserId:           "owner_user_id",
	Status:                "status",
	EntityNo:              "entity_no",
	ShopNo:                "shop_no",
	EntityName:            "entity_name",
	ShopName:              "shop_name",
	ShopNameNorm:          "shop_name_norm",
	EntityDraftJson:       "entity_draft_json",
	ShopDraftJson:         "shop_draft_json",
	EntitySubmittedJson:   "entity_submitted_json",
	ShopSubmittedJson:     "shop_submitted_json",
	LatestRejectJson:      "latest_reject_json",
	SubmittedAt:           "submitted_at",
	ReviewStartedAt:       "review_started_at",
	ReviewedAt:            "reviewed_at",
	ReviewerId:            "reviewer_id",
	ReviewComment:         "review_comment",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
}

// NewSellerApplicationDao creates and returns a new DAO object for table data access.
func NewSellerApplicationDao(handlers ...gdb.ModelHandler) *SellerApplicationDao {
	return &SellerApplicationDao{
		group:    "default",
		table:    "seller_application",
		columns:  sellerApplicationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SellerApplicationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SellerApplicationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SellerApplicationDao) Columns() SellerApplicationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SellerApplicationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SellerApplicationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SellerApplicationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
