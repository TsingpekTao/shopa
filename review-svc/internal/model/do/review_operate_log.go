// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewOperateLog is the golang structure of table review_operate_log for DAO operations like Where/Data.
type ReviewOperateLog struct {
	g.Meta       `orm:"table:review_operate_log, do:true"`
	Id           any         //
	ReviewNo     any         //
	OperatorType any         //
	OperatorId   any         //
	ActionCode   any         //
	ReasonCode   any         //
	Remark       any         //
	SnapshotJson any         //
	CreatedAt    *gtime.Time //
}
