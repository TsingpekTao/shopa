// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AfterSaleCaseDao is the data access object for the table after_sale_case.
type AfterSaleCaseDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  AfterSaleCaseColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// AfterSaleCaseColumns defines and stores column names for the table after_sale_case.
type AfterSaleCaseColumns struct {
	Id                   string //
	AfterSaleNo          string //
	OrderNo              string //
	SubOrderNo           string //
	ItemNo               string //
	UserId               string //
	ShopNo               string //
	SpuNo                string //
	SkuNo                string //
	Qty                  string //
	AfterSaleType        string //
	AfterSaleStatus      string //
	ApplyRefundAmount    string //
	ApprovedRefundAmount string //
	ReasonCode           string //
	ReasonDesc           string //
	EvidenceAssetIdsJson string //
	BuyerRemark          string //
	SellerReply          string //
	RejectReasonCode     string //
	CancelReasonCode     string //
	Version              string //
	ClosedAt             string //
	CreatedAt            string //
	UpdatedAt            string //
	DeletedAt            string //
}

// afterSaleCaseColumns holds the columns for the table after_sale_case.
var afterSaleCaseColumns = AfterSaleCaseColumns{
	Id:                   "id",
	AfterSaleNo:          "after_sale_no",
	OrderNo:              "order_no",
	SubOrderNo:           "sub_order_no",
	ItemNo:               "item_no",
	UserId:               "user_id",
	ShopNo:               "shop_no",
	SpuNo:                "spu_no",
	SkuNo:                "sku_no",
	Qty:                  "qty",
	AfterSaleType:        "after_sale_type",
	AfterSaleStatus:      "after_sale_status",
	ApplyRefundAmount:    "apply_refund_amount",
	ApprovedRefundAmount: "approved_refund_amount",
	ReasonCode:           "reason_code",
	ReasonDesc:           "reason_desc",
	EvidenceAssetIdsJson: "evidence_asset_ids_json",
	BuyerRemark:          "buyer_remark",
	SellerReply:          "seller_reply",
	RejectReasonCode:     "reject_reason_code",
	CancelReasonCode:     "cancel_reason_code",
	Version:              "version",
	ClosedAt:             "closed_at",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
	DeletedAt:            "deleted_at",
}

// NewAfterSaleCaseDao creates and returns a new DAO object for table data access.
func NewAfterSaleCaseDao(handlers ...gdb.ModelHandler) *AfterSaleCaseDao {
	return &AfterSaleCaseDao{
		group:    "default",
		table:    "after_sale_case",
		columns:  afterSaleCaseColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AfterSaleCaseDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AfterSaleCaseDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AfterSaleCaseDao) Columns() AfterSaleCaseColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AfterSaleCaseDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AfterSaleCaseDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AfterSaleCaseDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
