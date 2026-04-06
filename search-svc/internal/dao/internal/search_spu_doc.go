// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SearchSpuDocDao is the data access object for the table search_spu_doc.
type SearchSpuDocDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  SearchSpuDocColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// SearchSpuDocColumns defines and stores column names for the table search_spu_doc.
type SearchSpuDocColumns struct {
	Id                    string //
	SpuNo                 string //
	Title                 string //
	ShopNo                string //
	ShopName              string //
	CategoryNo            string //
	StoreCategoryId       string //
	StoreCategoryL1       string //
	StoreCategoryL2       string //
	StoreCategoryPathJson string //
	CoverAssetId          string //
	CoverUrl              string //
	MinPrice              string //
	MaxPrice              string //
	StockTotal            string //
	SalesCount            string //
	AvgScoreX100          string //
	ReviewTotal           string //
	ShopStatusCode        string //
	OnShelfStatusCode     string //
	AttrsJson             string //
	SourceVersion         string //
	SourceUpdatedAt       string //
	Deleted               string //
	DeletedAt             string //
	UpdatedAt             string //
	CreatedAt             string //
}

// searchSpuDocColumns holds the columns for the table search_spu_doc.
var searchSpuDocColumns = SearchSpuDocColumns{
	Id:                    "id",
	SpuNo:                 "spu_no",
	Title:                 "title",
	ShopNo:                "shop_no",
	ShopName:              "shop_name",
	CategoryNo:            "category_no",
	StoreCategoryId:       "store_category_id",
	StoreCategoryL1:       "store_category_l1",
	StoreCategoryL2:       "store_category_l2",
	StoreCategoryPathJson: "store_category_path_json",
	CoverAssetId:          "cover_asset_id",
	CoverUrl:              "cover_url",
	MinPrice:              "min_price",
	MaxPrice:              "max_price",
	StockTotal:            "stock_total",
	SalesCount:            "sales_count",
	AvgScoreX100:          "avg_score_x100",
	ReviewTotal:           "review_total",
	ShopStatusCode:        "shop_status_code",
	OnShelfStatusCode:     "on_shelf_status_code",
	AttrsJson:             "attrs_json",
	SourceVersion:         "source_version",
	SourceUpdatedAt:       "source_updated_at",
	Deleted:               "deleted",
	DeletedAt:             "deleted_at",
	UpdatedAt:             "updated_at",
	CreatedAt:             "created_at",
}

// NewSearchSpuDocDao creates and returns a new DAO object for table data access.
func NewSearchSpuDocDao(handlers ...gdb.ModelHandler) *SearchSpuDocDao {
	return &SearchSpuDocDao{
		group:    "default",
		table:    "search_spu_doc",
		columns:  searchSpuDocColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SearchSpuDocDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SearchSpuDocDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SearchSpuDocDao) Columns() SearchSpuDocColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SearchSpuDocDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SearchSpuDocDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SearchSpuDocDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
