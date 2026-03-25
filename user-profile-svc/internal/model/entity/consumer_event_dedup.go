// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ConsumerEventDedup is the golang structure for table consumer_event_dedup.
type ConsumerEventDedup struct {
	Id           uint64      `json:"id"           orm:"id"            description:"Primary key"`              // Primary key
	ConsumerName string      `json:"consumerName" orm:"consumer_name" description:"Consumer unique name"`     // Consumer unique name
	EventId      string      `json:"eventId"      orm:"event_id"      description:"Event id for idempotency"` // Event id for idempotency
	EventType    string      `json:"eventType"    orm:"event_type"    description:"Event type"`               // Event type
	ProcessedAt  *gtime.Time `json:"processedAt"  orm:"processed_at"  description:"Processed timestamp"`      // Processed timestamp
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"Created timestamp"`        // Created timestamp
}
