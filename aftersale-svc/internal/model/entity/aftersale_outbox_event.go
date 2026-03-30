// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AftersaleOutboxEvent is the golang structure for table aftersale_outbox_event.
type AftersaleOutboxEvent struct {
	Id            uint64      `json:"id"            orm:"id"             ` //
	EventId       string      `json:"eventId"       orm:"event_id"       ` //
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" ` //
	AggregateNo   string      `json:"aggregateNo"   orm:"aggregate_no"   ` //
	EventType     string      `json:"eventType"     orm:"event_type"     ` //
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   ` //
	Status        uint        `json:"status"        orm:"status"         ` //
	AvailableAt   *gtime.Time `json:"availableAt"   orm:"available_at"   ` //
	SentAt        *gtime.Time `json:"sentAt"        orm:"sent_at"        ` //
	RetryCount    uint        `json:"retryCount"    orm:"retry_count"    ` //
	LastError     string      `json:"lastError"     orm:"last_error"     ` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` //
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"     ` //
}
