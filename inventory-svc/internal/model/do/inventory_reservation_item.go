// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryReservationItem 是用于 DAO 操作（如 Where/Data）的 inventory_reservation_item 表 Go 结构体。
type InventoryReservationItem struct {
	g.Meta        `orm:"table:inventory_reservation_item, do:true"`
	Id            any         //
	ReservationNo any         //
	OrderNo       any         //
	SkuNo         any         //
	SpuNo         any         //
	ShopNo        any         //
	Qty           any         //
	StockVersion  any         // Version after reserve success
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
