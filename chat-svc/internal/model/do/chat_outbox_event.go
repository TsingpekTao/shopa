// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatOutboxEvent is the golang structure of table chat_outbox_event for DAO operations like Where/Data.
type ChatOutboxEvent struct {
	g.Meta       `orm:"table:chat_outbox_event, do:true"`
	Id           any         //
	EventId      any         //
	EventType    any         //
	AggregateNo  any         //
	PayloadJson  any         //
	Status       any         //
	RetryCount   any         //
	NextRetryAt  *gtime.Time //
	PublishedAt  *gtime.Time //
	ErrorMessage any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
