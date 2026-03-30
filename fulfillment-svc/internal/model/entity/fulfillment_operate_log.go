// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentOperateLog is the golang structure for table fulfillment_operate_log.
type FulfillmentOperateLog struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`                       //
	ShipmentNo   string      `json:"shipmentNo"   orm:"shipment_no"   description:""`                       //
	OperatorType string      `json:"operatorType" orm:"operator_type" description:"SYSTEM/SELLER/INTERNAL"` // SYSTEM/SELLER/INTERNAL
	OperatorId   string      `json:"operatorId"   orm:"operator_id"   description:""`                       //
	Action       string      `json:"action"       orm:"action"        description:""`                       //
	FromStatus   uint        `json:"fromStatus"   orm:"from_status"   description:""`                       //
	ToStatus     uint        `json:"toStatus"     orm:"to_status"     description:""`                       //
	Remark       string      `json:"remark"       orm:"remark"        description:""`                       //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`                       //
}
