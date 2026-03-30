// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatOutboxEvent is the golang structure for table chat_outbox_event.
type ChatOutboxEvent struct {
	Id           uint64      `json:"id"           orm:"id"            description:""` //
	EventId      string      `json:"eventId"      orm:"event_id"      description:""` //
	EventType    string      `json:"eventType"    orm:"event_type"    description:""` //
	AggregateNo  string      `json:"aggregateNo"  orm:"aggregate_no"  description:""` //
	PayloadJson  string      `json:"payloadJson"  orm:"payload_json"  description:""` //
	Status       int         `json:"status"       orm:"status"        description:""` //
	RetryCount   int         `json:"retryCount"   orm:"retry_count"   description:""` //
	NextRetryAt  *gtime.Time `json:"nextRetryAt"  orm:"next_retry_at" description:""` //
	PublishedAt  *gtime.Time `json:"publishedAt"  orm:"published_at"  description:""` //
	ErrorMessage string      `json:"errorMessage" orm:"error_message" description:""` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""` //
}
