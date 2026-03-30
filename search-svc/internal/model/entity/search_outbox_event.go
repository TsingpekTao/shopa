// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SearchOutboxEvent is the golang structure for table search_outbox_event.
type SearchOutboxEvent struct {
	Id            uint64      `json:"id"            orm:"id"             ` //
	EventId       string      `json:"eventId"       orm:"event_id"       ` //
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" ` //
	AggregateId   string      `json:"aggregateId"   orm:"aggregate_id"   ` //
	EventType     string      `json:"eventType"     orm:"event_type"     ` //
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   ` //
	Status        uint        `json:"status"        orm:"status"         ` //
	RetryCount    uint        `json:"retryCount"    orm:"retry_count"    ` //
	NextRetryAt   *gtime.Time `json:"nextRetryAt"   orm:"next_retry_at"  ` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` //
}
