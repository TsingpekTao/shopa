// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderOperateLog is the golang structure for table order_operate_log.
type OrderOperateLog struct {
	Id               uint64      `json:"id"               orm:"id"                 description:""` //
	OrderNo          string      `json:"orderNo"          orm:"order_no"           description:""` //
	SubOrderNo       string      `json:"subOrderNo"       orm:"sub_order_no"       description:""` //
	OperatorUserId   uint64      `json:"operatorUserId"   orm:"operator_user_id"   description:""` //
	OperatorTypeCode string      `json:"operatorTypeCode" orm:"operator_type_code" description:""` //
	ActionCode       string      `json:"actionCode"       orm:"action_code"        description:""` //
	BeforeStatus     string      `json:"beforeStatus"     orm:"before_status"      description:""` //
	AfterStatus      string      `json:"afterStatus"      orm:"after_status"       description:""` //
	DetailJson       string      `json:"detailJson"       orm:"detail_json"        description:""` //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:""` //
}
