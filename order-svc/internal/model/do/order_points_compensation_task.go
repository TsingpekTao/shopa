// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderPointsCompensationTask is the golang structure of table order_points_compensation_task for DAO operations like Where/Data.
type OrderPointsCompensationTask struct {
	g.Meta              `orm:"table:order_points_compensation_task, do:true"`
	Id                  any         //
	TaskNo              any         //
	OrderNo             any         //
	UserId              any         //
	PointsReservationNo any         //
	ActionCode          any         //
	TaskStatus          any         //
	RetryCount          any         //
	NextRetryAt         *gtime.Time //
	LastError           any         //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
