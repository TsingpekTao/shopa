// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsRefundAction is the golang structure for table points_refund_action.
type PointsRefundAction struct {
	Id                uint64      `json:"id"                orm:"id"                  ` //
	ActionNo          string      `json:"actionNo"          orm:"action_no"           ` //
	ActionType        string      `json:"actionType"        orm:"action_type"         ` // RETURN/REVERSE
	RefundNo          string      `json:"refundNo"          orm:"refund_no"           ` //
	OrderNo           string      `json:"orderNo"           orm:"order_no"            ` //
	SubOrderNo        string      `json:"subOrderNo"        orm:"sub_order_no"        ` //
	ShopNo            string      `json:"shopNo"            orm:"shop_no"             ` //
	UserId            uint64      `json:"userId"            orm:"user_id"             ` //
	IdempotencyKey    string      `json:"idempotencyKey"    orm:"idempotency_key"     ` //
	RequestedPoints   uint64      `json:"requestedPoints"   orm:"requested_points"    ` //
	EffectivePoints   uint64      `json:"effectivePoints"   orm:"effective_points"    ` //
	CashAmountCent    int64       `json:"cashAmountCent"    orm:"cash_amount_cent"    ` //
	ActionStatus      string      `json:"actionStatus"      orm:"action_status"       ` // SUCCESS
	ResultPayloadJson string      `json:"resultPayloadJson" orm:"result_payload_json" ` //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          ` //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          ` //
}
