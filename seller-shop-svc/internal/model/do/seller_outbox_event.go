// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerOutboxEvent is the golang structure of table seller_outbox_event for DAO operations like Where/Data.
type SellerOutboxEvent struct {
	g.Meta        `orm:"table:seller_outbox_event, do:true"`
	Id            any         //
	EventId       any         // Global unique event id
	EventType     any         // SellerRoleAssignRequested/...
	AggregateType any         // APPLICATION/SHOP/ENTITY
	AggregateNo   any         // application_no/shop_no/entity_no
	RequestId     any         // Trace request id
	PayloadJson   any         //
	Status        any         //
	AvailableAt   *gtime.Time //
	SentAt        *gtime.Time //
	FailCount     any         //
	LastError     any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
