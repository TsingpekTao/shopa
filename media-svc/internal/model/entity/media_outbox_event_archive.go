// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaOutboxEventArchive is the golang structure for table media_outbox_event_archive.
type MediaOutboxEventArchive struct {
	Id            uint64      `json:"id"            orm:"id"             description:""` //
	EventId       string      `json:"eventId"       orm:"event_id"       description:""` //
	EventType     string      `json:"eventType"     orm:"event_type"     description:""` //
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" description:""` //
	AggregateKey  string      `json:"aggregateKey"  orm:"aggregate_key"  description:""` //
	RequestId     string      `json:"requestId"     orm:"request_id"     description:""` //
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   description:""` //
	Status        uint        `json:"status"        orm:"status"         description:""` //
	AvailableAt   *gtime.Time `json:"availableAt"   orm:"available_at"   description:""` //
	SentAt        *gtime.Time `json:"sentAt"        orm:"sent_at"        description:""` //
	FailCount     uint        `json:"failCount"     orm:"fail_count"     description:""` //
	LastError     string      `json:"lastError"     orm:"last_error"     description:""` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""` //
	ArchivedAt    *gtime.Time `json:"archivedAt"    orm:"archived_at"    description:""` //
}
