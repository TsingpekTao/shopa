// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderInventoryLink is the golang structure of table order_inventory_link for DAO operations like Where/Data.
type OrderInventoryLink struct {
	g.Meta        `orm:"table:order_inventory_link, do:true"`
	Id            any         //
	OrderNo       any         //
	ReservationNo any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
