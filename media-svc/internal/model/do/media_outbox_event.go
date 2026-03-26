// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaOutboxEvent 是 media_outbox_event 表的 Go 结构体，供 DAO 的 Where/Data 等操作使用。
type MediaOutboxEvent struct {
	g.Meta        `orm:"table:media_outbox_event, do:true"`
	Id            any         //
	EventId       any         // 全局唯一事件 ID。
	EventType     any         // 事件类型，例如 AssetProcessingCompleted/AssetProcessingFailed/...。
	AggregateType any         // 聚合类型：ASSET/BINDING/SCENE。
	AggregateKey  any         // 聚合键：asset_id/biz_key/scene_code。
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
