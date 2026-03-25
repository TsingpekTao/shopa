// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamOutboxEventArchive is the golang structure for table iam_outbox_event_archive.
type IamOutboxEventArchive struct {
	Id          uint64      `json:"id"          orm:"id"            description:""`                                                   //
	EventId     string      `json:"eventId"     orm:"event_id"      description:""`                                                   //
	EventType   string      `json:"eventType"   orm:"event_type"    description:""`                                                   //
	PayloadJson string      `json:"payloadJson" orm:"payload_json"  description:""`                                                   //
	Status      uint        `json:"status"      orm:"status"        description:"Archive rows are usually SENT/FAILED/DLQ snapshots"` // Archive rows are usually SENT/FAILED/DLQ snapshots
	AvailableAt *gtime.Time `json:"availableAt" orm:"available_at"  description:""`                                                   //
	SentAt      *gtime.Time `json:"sentAt"      orm:"sent_at"       description:""`                                                   //
	FailCount   uint        `json:"failCount"   orm:"fail_count"    description:""`                                                   //
	LastError   string      `json:"lastError"   orm:"last_error"    description:""`                                                   //
	RetryCount  uint        `json:"retryCount"  orm:"retry_count"   description:""`                                                   //
	NextRetryAt *gtime.Time `json:"nextRetryAt" orm:"next_retry_at" description:""`                                                   //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:""`                                                   //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:""`                                                   //
	ArchivedAt  *gtime.Time `json:"archivedAt"  orm:"archived_at"   description:""`                                                   //
}
