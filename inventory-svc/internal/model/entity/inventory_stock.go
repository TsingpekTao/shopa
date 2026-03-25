// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryStock is the golang structure for table inventory_stock.
type InventoryStock struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	SkuNo        string      `json:"skuNo"        orm:"sku_no"        ` //
	SpuNo        string      `json:"spuNo"        orm:"spu_no"        ` //
	ShopNo       string      `json:"shopNo"       orm:"shop_no"       ` //
	TotalQty     uint64      `json:"totalQty"     orm:"total_qty"     ` //
	LockedQty    uint64      `json:"lockedQty"    orm:"locked_qty"    ` //
	AvailableQty uint64      `json:"availableQty" orm:"available_qty" ` //
	StockVersion uint64      `json:"stockVersion" orm:"stock_version" ` // Monotonic version for projection events
	StockStatus  uint        `json:"stockStatus"  orm:"stock_status"  ` //
	IsHot        int         `json:"isHot"        orm:"is_hot"        ` //
	RowVersion   uint64      `json:"rowVersion"   orm:"row_version"   ` // Internal optimistic lock version
	LastEventId  string      `json:"lastEventId"  orm:"last_event_id" ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
