// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderPayment is the golang structure of table order_payment for DAO operations like Where/Data.
type OrderPayment struct {
	g.Meta         `orm:"table:order_payment, do:true"`
	Id             any         //
	OrderNo        any         //
	PayNo          any         //
	PaymentEventId any         //
	PayChannel     any         //
	PayStatusCode  any         //
	ChannelTradeNo any         //
	PaidAmount     any         //
	PaidAt         *gtime.Time //
	RawPayload     any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
