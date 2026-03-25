// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryConsumerEventDedup 是用于 DAO 操作（如 Where/Data）的 inventory_consumer_event_dedup 表 Go 结构体。
type InventoryConsumerEventDedup struct {
	g.Meta       `orm:"table:inventory_consumer_event_dedup, do:true"`
	Id           any         //
	ConsumerName any         //
	EventId      any         //
	PayloadHash  any         //
	FirstSeenAt  *gtime.Time //
}
