// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderMain is the golang structure of table order_main for DAO operations like Where/Data.
type OrderMain struct {
	g.Meta           `orm:"table:order_main, do:true"`
	Id               any         //
	OrderNo          any         //
	UserId           any         //
	OrderStatus      any         //
	PaymentStatus    any         //
	ReservationNo    any         //
	GoodsAmount      any         //
	FreightAmount    any         //
	DiscountAmount   any         //
	PayableAmount    any         //
	PaidAmount       any         //
	BuyerRemark      any         //
	CancelReasonCode any         //
	PayDeadlineAt    *gtime.Time //
	PaidAt           *gtime.Time //
	ClosedAt         *gtime.Time //
	Version          any         //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
	DeletedAt        *gtime.Time //
}
