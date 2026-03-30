// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CartItemBackupDao is the data access object for the table cart_item_backup.
type CartItemBackupDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  CartItemBackupColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// CartItemBackupColumns defines and stores column names for the table cart_item_backup.
type CartItemBackupColumns struct {
	Id                string //
	UserId            string //
	SkuNo             string //
	SpuNo             string //
	ShopNo            string //
	Qty               string //
	Checked           string //
	Status            string // 1=ACTIVE,2=INVALID
	InvalidReasonCode string //
	SpuTitle          string //
	SkuName           string //
	SkuImageAssetId   string //
	SalePrice         string //
	MarketPrice       string //
	SaleAttrsJson     string //
	CreatedAt         string //
	UpdatedAt         string //
}

// cartItemBackupColumns holds the columns for the table cart_item_backup.
var cartItemBackupColumns = CartItemBackupColumns{
	Id:                "id",
	UserId:            "user_id",
	SkuNo:             "sku_no",
	SpuNo:             "spu_no",
	ShopNo:            "shop_no",
	Qty:               "qty",
	Checked:           "checked",
	Status:            "status",
	InvalidReasonCode: "invalid_reason_code",
	SpuTitle:          "spu_title",
	SkuName:           "sku_name",
	SkuImageAssetId:   "sku_image_asset_id",
	SalePrice:         "sale_price",
	MarketPrice:       "market_price",
	SaleAttrsJson:     "sale_attrs_json",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewCartItemBackupDao creates and returns a new DAO object for table data access.
func NewCartItemBackupDao(handlers ...gdb.ModelHandler) *CartItemBackupDao {
	return &CartItemBackupDao{
		group:    "default",
		table:    "cart_item_backup",
		columns:  cartItemBackupColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CartItemBackupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CartItemBackupDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CartItemBackupDao) Columns() CartItemBackupColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CartItemBackupDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CartItemBackupDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CartItemBackupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
