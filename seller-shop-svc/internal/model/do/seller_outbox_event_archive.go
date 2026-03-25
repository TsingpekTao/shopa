// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerOutboxEventArchive is the golang structure of table seller_outbox_event_archive for DAO operations like Where/Data.
type SellerOutboxEventArchive struct {
	g.Meta        `orm:"table:seller_outbox_event_archive, do:true"`
	Id            any         //
	EventId       any         //
	EventType     any         //
	AggregateType any         //
	AggregateNo   any         //
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
