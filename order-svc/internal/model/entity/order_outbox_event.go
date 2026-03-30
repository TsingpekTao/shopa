// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderOutboxEvent is the golang structure for table order_outbox_event.
type OrderOutboxEvent struct {
	Id            uint64      `json:"id"            orm:"id"             description:""` //
	EventId       string      `json:"eventId"       orm:"event_id"       description:""` //
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" description:""` //
	AggregateNo   string      `json:"aggregateNo"   orm:"aggregate_no"   description:""` //
	EventType     string      `json:"eventType"     orm:"event_type"     description:""` //
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   description:""` //
	Status        uint        `json:"status"        orm:"status"         description:""` //
	AvailableAt   *gtime.Time `json:"availableAt"   orm:"available_at"   description:""` //
	SentAt        *gtime.Time `json:"sentAt"        orm:"sent_at"        description:""` //
	RetryCount    uint        `json:"retryCount"    orm:"retry_count"    description:""` //
	LastError     string      `json:"lastError"     orm:"last_error"     description:""` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""` //
}
