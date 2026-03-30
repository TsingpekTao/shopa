// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewOperateLog is the golang structure for table review_operate_log.
type ReviewOperateLog struct {
	Id           uint64      `json:"id"           orm:"id"            description:""` //
	ReviewNo     string      `json:"reviewNo"     orm:"review_no"     description:""` //
	OperatorType string      `json:"operatorType" orm:"operator_type" description:""` //
	OperatorId   uint64      `json:"operatorId"   orm:"operator_id"   description:""` //
	ActionCode   string      `json:"actionCode"   orm:"action_code"   description:""` //
	ReasonCode   string      `json:"reasonCode"   orm:"reason_code"   description:""` //
	Remark       string      `json:"remark"       orm:"remark"        description:""` //
	SnapshotJson string      `json:"snapshotJson" orm:"snapshot_json" description:""` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""` //
}
