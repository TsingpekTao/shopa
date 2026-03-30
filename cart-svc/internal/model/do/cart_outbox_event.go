// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CartOutboxEvent is the golang structure of table cart_outbox_event for DAO operations like Where/Data.
type CartOutboxEvent struct {
	g.Meta        `orm:"table:cart_outbox_event, do:true"`
	Id            any         //
	EventId       any         //
	EventType     any         //
	AggregateType any         //
	AggregateId   any         //
	PayloadJson   any         //
	HeadersJson   any         //
	Status        any         // 0=NEW,1=PROCESSING,2=SENT,3=FAILED,4=DLQ
	AvailableAt   *gtime.Time //
	RetryCount    any         //
	LastError     any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
