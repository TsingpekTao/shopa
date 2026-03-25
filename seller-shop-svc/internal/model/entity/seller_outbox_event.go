// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerOutboxEvent is the golang structure for table seller_outbox_event.
type SellerOutboxEvent struct {
	Id            uint64      `json:"id"            orm:"id"             ` //
	EventId       string      `json:"eventId"       orm:"event_id"       ` // Global unique event id
	EventType     string      `json:"eventType"     orm:"event_type"     ` // SellerRoleAssignRequested/...
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" ` // APPLICATION/SHOP/ENTITY
	AggregateNo   string      `json:"aggregateNo"   orm:"aggregate_no"   ` // application_no/shop_no/entity_no
	RequestId     string      `json:"requestId"     orm:"request_id"     ` // Trace request id
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   ` //
	Status        uint        `json:"status"        orm:"status"         ` //
	AvailableAt   *gtime.Time `json:"availableAt"   orm:"available_at"   ` //
	SentAt        *gtime.Time `json:"sentAt"        orm:"sent_at"        ` //
	FailCount     uint        `json:"failCount"     orm:"fail_count"     ` //
	LastError     string      `json:"lastError"     orm:"last_error"     ` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` //
}
