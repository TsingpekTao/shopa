// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaOutboxEvent is the golang structure of table media_outbox_event for DAO operations like Where/Data.
type MediaOutboxEvent struct {
	g.Meta        `orm:"table:media_outbox_event, do:true"`
	Id            any         //
	EventId       any         // Global unique event id
	EventType     any         // AssetProcessingCompleted/AssetProcessingFailed/...
	AggregateType any         // ASSET/BINDING/SCENE
	AggregateKey  any         // asset_id/biz_key/scene_code
	RequestId     any         //
	PayloadJson   any         //
	Status        any         //
	AvailableAt   *gtime.Time //
	SentAt        *gtime.Time //
	FailCount     any         //
	LastError     any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
