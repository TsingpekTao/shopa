// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerApplicationOpLog is the golang structure for table seller_application_op_log.
type SellerApplicationOpLog struct {
	Id             uint64      `json:"id"             orm:"id"               ` //
	ApplicationNo  string      `json:"applicationNo"  orm:"application_no"   ` //
	OwnerUserId    uint64      `json:"ownerUserId"    orm:"owner_user_id"    ` //
	OperatorUserId uint64      `json:"operatorUserId" orm:"operator_user_id" ` // Admin or seller operator
	ActionCode     string      `json:"actionCode"     orm:"action_code"      ` // CREATE_DRAFT/SUBMIT/APPROVE/REJECT/RESUBMIT
	FromStatus     uint        `json:"fromStatus"     orm:"from_status"      ` //
	ToStatus       uint        `json:"toStatus"       orm:"to_status"        ` //
	Comment        string      `json:"comment"        orm:"comment"          ` //
	ExtraJson      string      `json:"extraJson"      orm:"extra_json"       ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       ` //
}
