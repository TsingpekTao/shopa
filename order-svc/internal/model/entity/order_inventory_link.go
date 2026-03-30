// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderInventoryLink is the golang structure for table order_inventory_link.
type OrderInventoryLink struct {
	Id            uint64      `json:"id"            orm:"id"             description:""` //
	OrderNo       string      `json:"orderNo"       orm:"order_no"       description:""` //
	ReservationNo string      `json:"reservationNo" orm:"reservation_no" description:""` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""` //
}
