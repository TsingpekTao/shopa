// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaConsumerEventDedup 是 media_consumer_event_dedup 表的 Go 结构体，供 DAO 的 Where/Data 等操作使用。
type MediaConsumerEventDedup struct {
	g.Meta       `orm:"table:media_consumer_event_dedup, do:true"`
	Id           any         //
	ConsumerName any         //
	EventId      any         //
	PayloadHash  any         //
	FirstSeenAt  *gtime.Time //
}
