// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationOutboxEvent is the golang structure for table notification_outbox_event.
type NotificationOutboxEvent struct {
	Id          uint64      `json:"id"          orm:"id"            ` //
	EventId     string      `json:"eventId"     orm:"event_id"      ` //
	Topic       string      `json:"topic"       orm:"topic"         ` //
	EventKey    string      `json:"eventKey"    orm:"event_key"     ` //
	PayloadJson string      `json:"payloadJson" orm:"payload_json"  ` //
	Status      uint        `json:"status"      orm:"status"        ` //
	NextRetryAt *gtime.Time `json:"nextRetryAt" orm:"next_retry_at" ` //
	RetryCount  uint        `json:"retryCount"  orm:"retry_count"   ` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    ` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    ` //
}
