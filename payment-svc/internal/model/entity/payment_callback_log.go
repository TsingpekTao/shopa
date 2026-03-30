// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentCallbackLog is the golang structure for table payment_callback_log.
type PaymentCallbackLog struct {
	Id                uint64      `json:"id"                orm:"id"                  ` //
	CallbackEventId   string      `json:"callbackEventId"   orm:"callback_event_id"   ` //
	PaymentNo         string      `json:"paymentNo"         orm:"payment_no"          ` //
	OrderNo           string      `json:"orderNo"           orm:"order_no"            ` //
	PayChannel        uint        `json:"payChannel"        orm:"pay_channel"         ` //
	GatewayStatusCode string      `json:"gatewayStatusCode" orm:"gateway_status_code" ` //
	PaidAmount        uint64      `json:"paidAmount"        orm:"paid_amount"         ` //
	RawPayload        string      `json:"rawPayload"        orm:"raw_payload"         ` //
	Signature         string      `json:"signature"         orm:"signature"           ` //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          ` //
}
