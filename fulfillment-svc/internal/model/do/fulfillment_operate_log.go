// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentOperateLog is the golang structure of table fulfillment_operate_log for DAO operations like Where/Data.
type FulfillmentOperateLog struct {
	g.Meta       `orm:"table:fulfillment_operate_log, do:true"`
	Id           any         //
	ShipmentNo   any         //
	OperatorType any         // SYSTEM/SELLER/INTERNAL
	OperatorId   any         //
	Action       any         //
	FromStatus   any         //
	ToStatus     any         //
	Remark       any         //
	CreatedAt    *gtime.Time //
}
