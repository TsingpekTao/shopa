// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderPayment is the golang structure for table order_payment.
type OrderPayment struct {
	Id             uint64      `json:"id"             orm:"id"               ` //
	OrderNo        string      `json:"orderNo"        orm:"order_no"         ` //
	PayNo          string      `json:"payNo"          orm:"pay_no"           ` //
	PaymentEventId string      `json:"paymentEventId" orm:"payment_event_id" ` //
	PayChannel     uint        `json:"payChannel"     orm:"pay_channel"      ` //
	PayStatusCode  string      `json:"payStatusCode"  orm:"pay_status_code"  ` //
	ChannelTradeNo string      `json:"channelTradeNo" orm:"channel_trade_no" ` //
	PaidAmount     uint64      `json:"paidAmount"     orm:"paid_amount"      ` //
	PaidAt         *gtime.Time `json:"paidAt"         orm:"paid_at"          ` //
	RawPayload     string      `json:"rawPayload"     orm:"raw_payload"      ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       ` //
}
