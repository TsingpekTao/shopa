// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewRecordDao is the data access object for the table review_record.
type ReviewRecordDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  ReviewRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// ReviewRecordColumns defines and stores column names for the table review_record.
type ReviewRecordColumns struct {
	Id               string //
	ReviewNo         string //
	OrderNo          string //
	SubOrderNo       string //
	ItemNo           string //
	UserId           string //
	ShopNo           string //
	SpuNo            string //
	SkuNo            string //
	Score            string //
	Content          string //
	MediasJson       string //
	Anonymous        string //
	AppendContent    string //
	AppendMediasJson string //
	AppendAt         string //
	SellerReply      string //
	SellerReplyAt    string //
	ReviewStatus     string //
	LikeCount        string //
	Version          string //
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// reviewRecordColumns holds the columns for the table review_record.
var reviewRecordColumns = ReviewRecordColumns{
	Id:               "id",
	ReviewNo:         "review_no",
	OrderNo:          "order_no",
	SubOrderNo:       "sub_order_no",
	ItemNo:           "item_no",
	UserId:           "user_id",
	ShopNo:           "shop_no",
	SpuNo:            "spu_no",
	SkuNo:            "sku_no",
	Score:            "score",
	Content:          "content",
	MediasJson:       "medias_json",
	Anonymous:        "anonymous",
	AppendContent:    "append_content",
	AppendMediasJson: "append_medias_json",
	AppendAt:         "append_at",
	SellerReply:      "seller_reply",
	SellerReplyAt:    "seller_reply_at",
	ReviewStatus:     "review_status",
	LikeCount:        "like_count",
	Version:          "version",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewReviewRecordDao creates and returns a new DAO object for table data access.
func NewReviewRecordDao(handlers ...gdb.ModelHandler) *ReviewRecordDao {
	return &ReviewRecordDao{
		group:    "default",
		table:    "review_record",
		columns:  reviewRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReviewRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReviewRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReviewRecordDao) Columns() ReviewRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReviewRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReviewRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReviewRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
