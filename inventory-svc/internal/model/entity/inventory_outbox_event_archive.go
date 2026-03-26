// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryOutboxEventArchive 是 inventory_outbox_event_archive 表的 Go 结构体。
type InventoryOutboxEventArchive struct {
	Id            uint64      `json:"id"            orm:"id"             ` //
	EventId       string      `json:"eventId"       orm:"event_id"       ` //
	EventType     string      `json:"eventType"     orm:"event_type"     ` //
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" ` //
	AggregateKey  string      `json:"aggregateKey"  orm:"aggregate_key"  ` //
	RequestId     string      `json:"requestId"     orm:"request_id"     ` //
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   ` //
	Status        uint        `json:"status"        orm:"status"         ` //
	AvailableAt   *gtime.Time `json:"availableAt"   orm:"available_at"   ` //
	SentAt        *gtime.Time `json:"sentAt"        orm:"sent_at"        ` //
	FailCount     uint        `json:"failCount"     orm:"fail_count"     ` //
	LastError     string      `json:"lastError"     orm:"last_error"     ` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` //
	ArchivedAt    *gtime.Time `json:"archivedAt"    orm:"archived_at"    ` //
}
