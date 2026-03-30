// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentOutboxEvent is the golang structure of table fulfillment_outbox_event for DAO operations like Where/Data.
type FulfillmentOutboxEvent struct {
	g.Meta        `orm:"table:fulfillment_outbox_event, do:true"`
	Id            any         //
	EventId       any         //
	EventType     any         //
	AggregateType any         //
	AggregateId   any         //
	PayloadJson   any         //
	Status        any         // 0 PENDING 1 PUBLISHED 2 FAILED
	RetryCount    any         //
	NextRetryAt   *gtime.Time //
	LastError     any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
