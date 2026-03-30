// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderSub is the golang structure of table order_sub for DAO operations like Where/Data.
type OrderSub struct {
	g.Meta         `orm:"table:order_sub, do:true"`
	Id             any         //
	SubOrderNo     any         //
	OrderNo        any         //
	ShopNo         any         //
	SubStatus      any         //
	GoodsAmount    any         //
	FreightAmount  any         //
	DiscountAmount any         //
	PayableAmount  any         //
	PaidAmount     any         //
	SellerRemark   any         //
	BuyerRemark    any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
