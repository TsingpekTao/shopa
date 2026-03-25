// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryOutboxEvent 是用于 DAO 操作（如 Where/Data）的 inventory_outbox_event 表 Go 结构体。
type InventoryOutboxEvent struct {
	g.Meta        `orm:"table:inventory_outbox_event, do:true"`
	Id            any         //
	EventId       any         //
	EventType     any         //
	AggregateType any         // STOCK/RESERVATION
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
}
