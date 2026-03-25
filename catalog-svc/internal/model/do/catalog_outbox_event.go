// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogOutboxEvent is the golang structure of table catalog_outbox_event for DAO operations like Where/Data.
type CatalogOutboxEvent struct {
	g.Meta        `orm:"table:catalog_outbox_event, do:true"`
	Id            any         //
	EventId       any         //
	EventType     any         //
	AggregateType any         // SPU/SKU/REVIEW
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
