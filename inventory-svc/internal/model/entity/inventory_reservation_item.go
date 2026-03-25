// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryReservationItem is the golang structure for table inventory_reservation_item.
type InventoryReservationItem struct {
	Id            uint64      `json:"id"            orm:"id"             ` //
	ReservationNo string      `json:"reservationNo" orm:"reservation_no" ` //
	OrderNo       string      `json:"orderNo"       orm:"order_no"       ` //
	SkuNo         string      `json:"skuNo"         orm:"sku_no"         ` //
	SpuNo         string      `json:"spuNo"         orm:"spu_no"         ` //
	ShopNo        string      `json:"shopNo"        orm:"shop_no"        ` //
	Qty           uint        `json:"qty"           orm:"qty"            ` //
	StockVersion  uint64      `json:"stockVersion"  orm:"stock_version"  ` // Version after reserve success
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` //
}
