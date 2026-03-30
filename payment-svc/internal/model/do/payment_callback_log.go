// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentCallbackLog is the golang structure of table payment_callback_log for DAO operations like Where/Data.
type PaymentCallbackLog struct {
	g.Meta            `orm:"table:payment_callback_log, do:true"`
	Id                any         //
	CallbackEventId   any         //
	PaymentNo         any         //
	OrderNo           any         //
	PayChannel        any         //
	GatewayStatusCode any         //
	PaidAmount        any         //
	RawPayload        any         //
	Signature         any         //
	CreatedAt         *gtime.Time //
}
