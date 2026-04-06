// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderInventoryLink is the golang structure for table order_inventory_link.
type OrderInventoryLink struct {
	Id            uint64      `json:"id"            orm:"id"             ` //
	OrderNo       string      `json:"orderNo"       orm:"order_no"       ` //
	ReservationNo string      `json:"reservationNo" orm:"reservation_no" ` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` //
}
