// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CartOutboxEvent is the golang structure for table cart_outbox_event.
type CartOutboxEvent struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`                                         //
	EventId       string      `json:"eventId"       orm:"event_id"       description:""`                                         //
	EventType     string      `json:"eventType"     orm:"event_type"     description:""`                                         //
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" description:""`                                         //
	AggregateId   string      `json:"aggregateId"   orm:"aggregate_id"   description:""`                                         //
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   description:""`                                         //
	HeadersJson   string      `json:"headersJson"   orm:"headers_json"   description:""`                                         //
	Status        uint        `json:"status"        orm:"status"         description:"0=NEW,1=PROCESSING,2=SENT,3=FAILED,4=DLQ"` // 0=NEW,1=PROCESSING,2=SENT,3=FAILED,4=DLQ
	AvailableAt   *gtime.Time `json:"availableAt"   orm:"available_at"   description:""`                                         //
	RetryCount    uint        `json:"retryCount"    orm:"retry_count"    description:""`                                         //
	LastError     string      `json:"lastError"     orm:"last_error"     description:""`                                         //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`                                         //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`                                         //
}
