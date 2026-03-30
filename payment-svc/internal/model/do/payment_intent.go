// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentIntent is the golang structure of table payment_intent for DAO operations like Where/Data.
type PaymentIntent struct {
	g.Meta          `orm:"table:payment_intent, do:true"`
	Id              any         //
	PaymentNo       any         //
	OrderNo         any         //
	UserId          any         //
	PayChannel      any         //
	Status          any         //
	PayableAmount   any         //
	RefundedAmount  any         //
	CurrencyCode    any         //
	OrderExpireAt   *gtime.Time //
	GatewayExpireAt *gtime.Time //
	ExternalTradeNo any         //
	PaidAt          *gtime.Time //
	ClosedAt        *gtime.Time //
	Version         any         //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       any         //
}
