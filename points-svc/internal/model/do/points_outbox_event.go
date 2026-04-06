// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsOutboxEvent is the golang structure of table points_outbox_event for DAO operations like Where/Data.
type PointsOutboxEvent struct {
	g.Meta            `orm:"table:points_outbox_event, do:true"`
	Id                any         //
	EventNo           any         //
	AggregateType     any         //
	AggregateNo       any         //
	EventType         any         //
	PayloadJson       any         //
	PublishStatusCode any         // PENDING/SENT/FAILED/DEAD
	RetryCount        any         //
	NextRetryAt       *gtime.Time //
	LastError         any         //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}
