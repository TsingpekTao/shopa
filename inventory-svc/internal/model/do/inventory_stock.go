// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryStock is the golang structure of table inventory_stock for DAO operations like Where/Data.
type InventoryStock struct {
	g.Meta       `orm:"table:inventory_stock, do:true"`
	Id           any         //
	SkuNo        any         //
	SpuNo        any         //
	ShopNo       any         //
	TotalQty     any         //
	LockedQty    any         //
	AvailableQty any         //
	StockVersion any         // Monotonic version for projection events
	StockStatus  any         //
	IsHot        any         //
	RowVersion   any         // Internal optimistic lock version
	LastEventId  any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
