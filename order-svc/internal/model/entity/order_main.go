// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderMain is the golang structure for table order_main.
type OrderMain struct {
	Id                       uint64      `json:"id"                       orm:"id"                          ` //
	OrderNo                  string      `json:"orderNo"                  orm:"order_no"                    ` //
	UserId                   uint64      `json:"userId"                   orm:"user_id"                     ` //
	OrderStatus              uint        `json:"orderStatus"              orm:"order_status"                ` //
	PaymentStatus            uint        `json:"paymentStatus"            orm:"payment_status"              ` //
	ReservationNo            string      `json:"reservationNo"            orm:"reservation_no"              ` //
	PointsReservationNo      string      `json:"pointsReservationNo"      orm:"points_reservation_no"       ` //
	GoodsAmount              uint64      `json:"goodsAmount"              orm:"goods_amount"                ` //
	FreightAmount            uint64      `json:"freightAmount"            orm:"freight_amount"              ` //
	DiscountAmount           uint64      `json:"discountAmount"           orm:"discount_amount"             ` //
	PayableAmount            uint64      `json:"payableAmount"            orm:"payable_amount"              ` //
	PaidAmount               uint64      `json:"paidAmount"               orm:"paid_amount"                 ` //
	PointsUsed               uint64      `json:"pointsUsed"               orm:"points_used"                 ` //
	PointsDiscountAmount     uint64      `json:"pointsDiscountAmount"     orm:"points_discount_amount"      ` //
	PointsRuleSnapshotJson   string      `json:"pointsRuleSnapshotJson"   orm:"points_rule_snapshot_json"   ` //
	PointsRuleSnapshotDigest string      `json:"pointsRuleSnapshotDigest" orm:"points_rule_snapshot_digest" ` //
	BuyerRemark              string      `json:"buyerRemark"              orm:"buyer_remark"                ` //
	CancelReasonCode         uint        `json:"cancelReasonCode"         orm:"cancel_reason_code"          ` //
	PayDeadlineAt            *gtime.Time `json:"payDeadlineAt"            orm:"pay_deadline_at"             ` //
	PaidAt                   *gtime.Time `json:"paidAt"                   orm:"paid_at"                     ` //
	ClosedAt                 *gtime.Time `json:"closedAt"                 orm:"closed_at"                   ` //
	Version                  uint64      `json:"version"                  orm:"version"                     ` //
	CreatedAt                *gtime.Time `json:"createdAt"                orm:"created_at"                  ` //
	UpdatedAt                *gtime.Time `json:"updatedAt"                orm:"updated_at"                  ` //
	DeletedAt                *gtime.Time `json:"deletedAt"                orm:"deleted_at"                  ` //
}
