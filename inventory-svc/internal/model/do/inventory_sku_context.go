// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventorySkuContext is the golang structure of table inventory_sku_context for DAO operations like Where/Data.
type InventorySkuContext struct {
	g.Meta    `orm:"table:inventory_sku_context, do:true"`
	Id        any         //
	SkuNo     any         //
	SpuNo     any         //
	ShopNo    any         //
	Enabled   any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
