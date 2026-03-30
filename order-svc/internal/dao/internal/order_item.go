// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderItemDao is the data access object for the table order_item.
type OrderItemDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  OrderItemColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// OrderItemColumns defines and stores column names for the table order_item.
type OrderItemColumns struct {
	Id              string //
	ItemNo          string //
	OrderNo         string //
	SubOrderNo      string //
	ShopNo          string //
	SpuNo           string //
	SkuNo           string //
	SpuTitle        string //
	SkuName         string //
	SkuImageAssetId string //
	Qty             string //
	SalePrice       string //
	MarketPrice     string //
	SaleAttrsJson   string //
	CreatedAt       string //
	UpdatedAt       string //
}

// orderItemColumns holds the columns for the table order_item.
var orderItemColumns = OrderItemColumns{
	Id:              "id",
	ItemNo:          "item_no",
	OrderNo:         "order_no",
	SubOrderNo:      "sub_order_no",
	ShopNo:          "shop_no",
	SpuNo:           "spu_no",
	SkuNo:           "sku_no",
	SpuTitle:        "spu_title",
	SkuName:         "sku_name",
	SkuImageAssetId: "sku_image_asset_id",
	Qty:             "qty",
	SalePrice:       "sale_price",
	MarketPrice:     "market_price",
	SaleAttrsJson:   "sale_attrs_json",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewOrderItemDao creates and returns a new DAO object for table data access.
func NewOrderItemDao(handlers ...gdb.ModelHandler) *OrderItemDao {
	return &OrderItemDao{
		group:    "default",
		table:    "order_item",
		columns:  orderItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrderItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrderItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrderItemDao) Columns() OrderItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrderItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrderItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OrderItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
