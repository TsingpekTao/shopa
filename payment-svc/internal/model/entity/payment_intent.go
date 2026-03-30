// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentIntent is the golang structure for table payment_intent.
type PaymentIntent struct {
	Id              uint64      `json:"id"              orm:"id"                ` //
	PaymentNo       string      `json:"paymentNo"       orm:"payment_no"        ` //
	OrderNo         string      `json:"orderNo"         orm:"order_no"          ` //
	UserId          uint64      `json:"userId"          orm:"user_id"           ` //
	PayChannel      uint        `json:"payChannel"      orm:"pay_channel"       ` //
	Status          uint        `json:"status"          orm:"status"            ` //
	PayableAmount   uint64      `json:"payableAmount"   orm:"payable_amount"    ` //
	RefundedAmount  uint64      `json:"refundedAmount"  orm:"refunded_amount"   ` //
	CurrencyCode    string      `json:"currencyCode"    orm:"currency_code"     ` //
	OrderExpireAt   *gtime.Time `json:"orderExpireAt"   orm:"order_expire_at"   ` //
	GatewayExpireAt *gtime.Time `json:"gatewayExpireAt" orm:"gateway_expire_at" ` //
	ExternalTradeNo string      `json:"externalTradeNo" orm:"external_trade_no" ` //
	PaidAt          *gtime.Time `json:"paidAt"          orm:"paid_at"           ` //
	ClosedAt        *gtime.Time `json:"closedAt"        orm:"closed_at"         ` //
	Version         uint64      `json:"version"         orm:"version"           ` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"        ` //
	DeletedAt       int64       `json:"deletedAt"       orm:"deleted_at"        ` //
}
