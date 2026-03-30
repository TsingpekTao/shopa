// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItem is the golang structure of table order_item for DAO operations like Where/Data.
type OrderItem struct {
	g.Meta          `orm:"table:order_item, do:true"`
	Id              any         //
	ItemNo          any         //
	OrderNo         any         //
	SubOrderNo      any         //
	ShopNo          any         //
	SpuNo           any         //
	SkuNo           any         //
	SpuTitle        any         //
	SkuName         any         //
	SkuImageAssetId any         //
	Qty             any         //
	SalePrice       any         //
	MarketPrice     any         //
	SaleAttrsJson   any         //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
