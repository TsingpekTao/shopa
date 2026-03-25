// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaOutboxEvent is the golang structure for table media_outbox_event.
type MediaOutboxEvent struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`                                                   //
	EventId       string      `json:"eventId"       orm:"event_id"       description:"Global unique event id"`                             // Global unique event id
	EventType     string      `json:"eventType"     orm:"event_type"     description:"AssetProcessingCompleted/AssetProcessingFailed/..."` // AssetProcessingCompleted/AssetProcessingFailed/...
	AggregateType string      `json:"aggregateType" orm:"aggregate_type" description:"ASSET/BINDING/SCENE"`                                // ASSET/BINDING/SCENE
	AggregateKey  string      `json:"aggregateKey"  orm:"aggregate_key"  description:"asset_id/biz_key/scene_code"`                        // asset_id/biz_key/scene_code
	RequestId     string      `json:"requestId"     orm:"request_id"     description:""`                                                   //
	PayloadJson   string      `json:"payloadJson"   orm:"payload_json"   description:""`                                                   //
	Status        uint        `json:"status"        orm:"status"         description:""`                                                   //
	AvailableAt   *gtime.Time `json:"availableAt"   orm:"available_at"   description:""`                                                   //
	SentAt        *gtime.Time `json:"sentAt"        orm:"sent_at"        description:""`                                                   //
	FailCount     uint        `json:"failCount"     orm:"fail_count"     description:""`                                                   //
	LastError     string      `json:"lastError"     orm:"last_error"     description:""`                                                   //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`                                                   //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`                                                   //
}
