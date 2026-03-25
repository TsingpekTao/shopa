// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryOutboxEvent is the golang structure of table inventory_outbox_event for DAO operations like Where/Data.
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
