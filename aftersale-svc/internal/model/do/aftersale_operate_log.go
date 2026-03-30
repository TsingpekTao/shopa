// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AftersaleOperateLog is the golang structure of table aftersale_operate_log for DAO operations like Where/Data.
type AftersaleOperateLog struct {
	g.Meta       `orm:"table:aftersale_operate_log, do:true"`
	Id           any         //
	AfterSaleNo  any         //
	OperatorType any         //
	OperatorId   any         //
	ActionCode   any         //
	FromStatus   any         //
	ToStatus     any         //
	ReasonCode   any         //
	Remark       any         //
	CreatedAt    *gtime.Time //
}
