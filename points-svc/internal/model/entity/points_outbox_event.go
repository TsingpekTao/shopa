// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsOutboxEvent is the golang structure for table points_outbox_event.
type PointsOutboxEvent struct {
	Id                uint64      `json:"id"                orm:"id"                  ` //
	EventNo           string      `json:"eventNo"           orm:"event_no"            ` //
	AggregateType     string      `json:"aggregateType"     orm:"aggregate_type"      ` //
	AggregateNo       string      `json:"aggregateNo"       orm:"aggregate_no"        ` //
	EventType         string      `json:"eventType"         orm:"event_type"          ` //
	PayloadJson       string      `json:"payloadJson"       orm:"payload_json"        ` //
	PublishStatusCode string      `json:"publishStatusCode" orm:"publish_status_code" ` // PENDING/SENT/FAILED/DEAD
	RetryCount        uint        `json:"retryCount"        orm:"retry_count"         ` //
	NextRetryAt       *gtime.Time `json:"nextRetryAt"       orm:"next_retry_at"       ` //
	LastError         string      `json:"lastError"         orm:"last_error"          ` //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          ` //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          ` //
}
