// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AftersaleOperateLog is the golang structure for table aftersale_operate_log.
type AftersaleOperateLog struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	AfterSaleNo  string      `json:"afterSaleNo"  orm:"after_sale_no" ` //
	OperatorType string      `json:"operatorType" orm:"operator_type" ` //
	OperatorId   uint64      `json:"operatorId"   orm:"operator_id"   ` //
	ActionCode   string      `json:"actionCode"   orm:"action_code"   ` //
	FromStatus   uint        `json:"fromStatus"   orm:"from_status"   ` //
	ToStatus     uint        `json:"toStatus"     orm:"to_status"     ` //
	ReasonCode   string      `json:"reasonCode"   orm:"reason_code"   ` //
	Remark       string      `json:"remark"       orm:"remark"        ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
}
