// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SearchOutboxEvent is the golang structure of table search_outbox_event for DAO operations like Where/Data.
type SearchOutboxEvent struct {
	g.Meta        `orm:"table:search_outbox_event, do:true"`
	Id            any         //
	EventId       any         //
	AggregateType any         //
	AggregateId   any         //
	EventType     any         //
	PayloadJson   any         //
	Status        any         //
	RetryCount    any         //
	NextRetryAt   *gtime.Time //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
