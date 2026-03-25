// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamOutboxEvent is the golang structure for table iam_outbox_event.
type IamOutboxEvent struct {
	Id          uint64      `json:"id"          orm:"id"            description:""`                                         //
	EventId     string      `json:"eventId"     orm:"event_id"      description:"Global unique event id"`                   // Global unique event id
	EventType   string      `json:"eventType"   orm:"event_type"    description:""`                                         //
	PayloadJson string      `json:"payloadJson" orm:"payload_json"  description:""`                                         //
	Status      uint        `json:"status"      orm:"status"        description:"1 NEW,2 PROCESSING,3 SENT,4 FAILED,5 DLQ"` // 1 NEW,2 PROCESSING,3 SENT,4 FAILED,5 DLQ
	AvailableAt *gtime.Time `json:"availableAt" orm:"available_at"  description:"Next dispatch schedule time"`              // Next dispatch schedule time
	SentAt      *gtime.Time `json:"sentAt"      orm:"sent_at"       description:"Successful dispatch time"`                 // Successful dispatch time
	FailCount   uint        `json:"failCount"   orm:"fail_count"    description:"Dispatch failure count"`                   // Dispatch failure count
	LastError   string      `json:"lastError"   orm:"last_error"    description:"Last dispatch error text"`                 // Last dispatch error text
	RetryCount  uint        `json:"retryCount"  orm:"retry_count"   description:""`                                         //
	NextRetryAt *gtime.Time `json:"nextRetryAt" orm:"next_retry_at" description:""`                                         //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:""`                                         //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:""`                                         //
}
