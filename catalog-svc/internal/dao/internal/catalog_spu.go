// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogSpuDao is the data access object for the table catalog_spu.
type CatalogSpuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CatalogSpuColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CatalogSpuColumns defines and stores column names for the table catalog_spu.
type CatalogSpuColumns struct {
	Id                      string //
	SpuNo                   string //
	ShopNo                  string //
	Title                   string //
	SubTitle                string //
	CategoryId              string //
	BrandNo                 string //
	MainImageAssetIdsJson   string //
	DetailImageAssetIdsJson string //
	SpuStatus               string //
	SpuStockStatus          string //
	MinSalePrice            string //
	MaxSalePrice            string //
	MinMarketPrice          string //
	MaxMarketPrice          string //
	PublishTime             string //
	Version                 string //
	ReviewStatus            string // 1 PENDING, 2 APPROVED, 3 REJECTED
	RejectReasonCode        string //
	RejectComment           string //
	ReviewComment           string //
	ReviewerId              string //
	SubmittedAt             string //
	ReviewedAt              string //
	SoldCount               string //
	DeletedAt               string //
	CreatedAt               string //
	UpdatedAt               string //
}

// catalogSpuColumns holds the columns for the table catalog_spu.
var catalogSpuColumns = CatalogSpuColumns{
	Id:                      "id",
	SpuNo:                   "spu_no",
	ShopNo:                  "shop_no",
	Title:                   "title",
	SubTitle:                "sub_title",
	CategoryId:              "category_id",
	BrandNo:                 "brand_no",
	MainImageAssetIdsJson:   "main_image_asset_ids_json",
	DetailImageAssetIdsJson: "detail_image_asset_ids_json",
	SpuStatus:               "spu_status",
	SpuStockStatus:          "spu_stock_status",
	MinSalePrice:            "min_sale_price",
	MaxSalePrice:            "max_sale_price",
	MinMarketPrice:          "min_market_price",
	MaxMarketPrice:          "max_market_price",
	PublishTime:             "publish_time",
	Version:                 "version",
	ReviewStatus:            "review_status",
	RejectReasonCode:        "reject_reason_code",
	RejectComment:           "reject_comment",
	ReviewComment:           "review_comment",
	ReviewerId:              "reviewer_id",
	SubmittedAt:             "submitted_at",
	ReviewedAt:              "reviewed_at",
	SoldCount:               "sold_count",
	DeletedAt:               "deleted_at",
	CreatedAt:               "created_at",
	UpdatedAt:               "updated_at",
}

// NewCatalogSpuDao creates and returns a new DAO object for table data access.
func NewCatalogSpuDao(handlers ...gdb.ModelHandler) *CatalogSpuDao {
	return &CatalogSpuDao{
		group:    "default",
		table:    "catalog_spu",
		columns:  catalogSpuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogSpuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogSpuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogSpuDao) Columns() CatalogSpuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogSpuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogSpuDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogSpuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
