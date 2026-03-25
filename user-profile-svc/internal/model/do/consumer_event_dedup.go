// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ConsumerEventDedup is the golang structure of table consumer_event_dedup for DAO operations like Where/Data.
type ConsumerEventDedup struct {
	g.Meta       `orm:"table:consumer_event_dedup, do:true"`
	Id           any         // Primary key
	ConsumerName any         // Consumer unique name
	EventId      any         // Event id for idempotency
	EventType    any         // Event type
	ProcessedAt  *gtime.Time // Processed timestamp
	CreatedAt    *gtime.Time // Created timestamp
}
