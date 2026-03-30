// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderOperateLog is the golang structure of table order_operate_log for DAO operations like Where/Data.
type OrderOperateLog struct {
	g.Meta           `orm:"table:order_operate_log, do:true"`
	Id               any         //
	OrderNo          any         //
	SubOrderNo       any         //
	OperatorUserId   any         //
	OperatorTypeCode any         //
	ActionCode       any         //
	BeforeStatus     any         //
	AfterStatus      any         //
	DetailJson       any         //
	CreatedAt        *gtime.Time //
}
