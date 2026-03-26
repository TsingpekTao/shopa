// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaOutboxEventArchive 是 media_outbox_event_archive 表的 Go 结构体，供 DAO 的 Where/Data 等操作使用。
type MediaOutboxEventArchive struct {
	g.Meta        `orm:"table:media_outbox_event_archive, do:true"`
	Id            any         //
	EventId       any         //
	EventType     any         //
	AggregateType any         //
	AggregateKey  any         //
	RequestId     any         //
	PayloadJson   any         //
	Status        any         //
	AvailableAt   *gtime.Time //
	SentAt        *gtime.Time //
	FailCount     any         //
	LastError     any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	ArchivedAt    *gtime.Time //
}
