// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AftersaleOutboxEvent is the golang structure of table aftersale_outbox_event for DAO operations like Where/Data.
type AftersaleOutboxEvent struct {
	g.Meta        `orm:"table:aftersale_outbox_event, do:true"`
	Id            any         //
	EventId       any         //
	AggregateType any         //
	AggregateNo   any         //
	EventType     any         //
	PayloadJson   any         //
	Status        any         //
	AvailableAt   *gtime.Time //
	SentAt        *gtime.Time //
	RetryCount    any         //
	LastError     any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
