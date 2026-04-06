// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsRefundAction is the golang structure of table points_refund_action for DAO operations like Where/Data.
type PointsRefundAction struct {
	g.Meta            `orm:"table:points_refund_action, do:true"`
	Id                any         //
	ActionNo          any         //
	ActionType        any         // RETURN/REVERSE
	RefundNo          any         //
	OrderNo           any         //
	SubOrderNo        any         //
	ShopNo            any         //
	UserId            any         //
	IdempotencyKey    any         //
	RequestedPoints   any         //
	EffectivePoints   any         //
	CashAmountCent    any         //
	ActionStatus      any         // SUCCESS
	ResultPayloadJson any         //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}
