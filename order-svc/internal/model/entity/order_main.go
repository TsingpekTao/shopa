// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderMain is the golang structure for table order_main.
type OrderMain struct {
	Id               uint64      `json:"id"               orm:"id"                 description:""` //
	OrderNo          string      `json:"orderNo"          orm:"order_no"           description:""` //
	UserId           uint64      `json:"userId"           orm:"user_id"            description:""` //
	OrderStatus      uint        `json:"orderStatus"      orm:"order_status"       description:""` //
	PaymentStatus    uint        `json:"paymentStatus"    orm:"payment_status"     description:""` //
	ReservationNo    string      `json:"reservationNo"    orm:"reservation_no"     description:""` //
	GoodsAmount      uint64      `json:"goodsAmount"      orm:"goods_amount"       description:""` //
	FreightAmount    uint64      `json:"freightAmount"    orm:"freight_amount"     description:""` //
	DiscountAmount   uint64      `json:"discountAmount"   orm:"discount_amount"    description:""` //
	PayableAmount    uint64      `json:"payableAmount"    orm:"payable_amount"     description:""` //
	PaidAmount       uint64      `json:"paidAmount"       orm:"paid_amount"        description:""` //
	BuyerRemark      string      `json:"buyerRemark"      orm:"buyer_remark"       description:""` //
	CancelReasonCode uint        `json:"cancelReasonCode" orm:"cancel_reason_code" description:""` //
	PayDeadlineAt    *gtime.Time `json:"payDeadlineAt"    orm:"pay_deadline_at"    description:""` //
	PaidAt           *gtime.Time `json:"paidAt"           orm:"paid_at"            description:""` //
	ClosedAt         *gtime.Time `json:"closedAt"         orm:"closed_at"          description:""` //
	Version          uint64      `json:"version"          orm:"version"            description:""` //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:""` //
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         description:""` //
	DeletedAt        *gtime.Time `json:"deletedAt"        orm:"deleted_at"         description:""` //
}
