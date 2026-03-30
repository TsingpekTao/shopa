// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentOutboxEvent is the golang structure of table payment_outbox_event for DAO operations like Where/Data.
type PaymentOutboxEvent struct {
	g.Meta      `orm:"table:payment_outbox_event, do:true"`
	Id          any         //
	EventId     any         //
	Topic       any         //
	EventKey    any         //
	PayloadJson any         //
	Status      any         //
	NextRetryAt *gtime.Time //
	RetryCount  any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
