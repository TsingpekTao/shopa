// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewOutboxEvent is the golang structure for table review_outbox_event.
type ReviewOutboxEvent struct {
	Id            uint64      `json:"id"            orm:"id"             description:""` //
	EventId       string      `json:"eventId"       orm:"event_id"       description:""` //
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" description:""` //
	AggregateId   string      `json:"aggregateId"   orm:"aggregate_id"   description:""` //
	EventType     string      `json:"eventType"     orm:"event_type"     description:""` //
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   description:""` //
	Status        uint        `json:"status"        orm:"status"         description:""` //
	RetryCount    uint        `json:"retryCount"    orm:"retry_count"    description:""` //
	NextRetryAt   *gtime.Time `json:"nextRetryAt"   orm:"next_retry_at"  description:""` //
	PublishedAt   *gtime.Time `json:"publishedAt"   orm:"published_at"   description:""` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""` //
}
