// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderPayment is the golang structure for table order_payment.
type OrderPayment struct {
	Id             uint64      `json:"id"             orm:"id"               description:""` //
	OrderNo        string      `json:"orderNo"        orm:"order_no"         description:""` //
	PayNo          string      `json:"payNo"          orm:"pay_no"           description:""` //
	PaymentEventId string      `json:"paymentEventId" orm:"payment_event_id" description:""` //
	PayChannel     uint        `json:"payChannel"     orm:"pay_channel"      description:""` //
	PayStatusCode  string      `json:"payStatusCode"  orm:"pay_status_code"  description:""` //
	ChannelTradeNo string      `json:"channelTradeNo" orm:"channel_trade_no" description:""` //
	PaidAmount     uint64      `json:"paidAmount"     orm:"paid_amount"      description:""` //
	PaidAt         *gtime.Time `json:"paidAt"         orm:"paid_at"          description:""` //
	RawPayload     string      `json:"rawPayload"     orm:"raw_payload"      description:""` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:""` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:""` //
}
