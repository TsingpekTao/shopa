// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerShopOpLog is the golang structure for table seller_shop_op_log.
type SellerShopOpLog struct {
	Id             uint64      `json:"id"             orm:"id"               ` //
	ShopNo         string      `json:"shopNo"         orm:"shop_no"          ` //
	OwnerUserId    uint64      `json:"ownerUserId"    orm:"owner_user_id"    ` //
	OperatorUserId uint64      `json:"operatorUserId" orm:"operator_user_id" ` //
	ActionCode     string      `json:"actionCode"     orm:"action_code"      ` // PROVISION_START/ACTIVATE/FREEZE/CLOSE/REOPEN
	FromStatus     uint        `json:"fromStatus"     orm:"from_status"      ` //
	ToStatus       uint        `json:"toStatus"       orm:"to_status"        ` //
	ReasonCode     string      `json:"reasonCode"     orm:"reason_code"      ` //
	Reason         string      `json:"reason"         orm:"reason"           ` //
	RequestId      string      `json:"requestId"      orm:"request_id"       ` //
	EventId        string      `json:"eventId"        orm:"event_id"         ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       ` //
}
