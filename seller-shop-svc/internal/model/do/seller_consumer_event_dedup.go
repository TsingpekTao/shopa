// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerConsumerEventDedup is the golang structure of table seller_consumer_event_dedup for DAO operations like Where/Data.
type SellerConsumerEventDedup struct {
	g.Meta       `orm:"table:seller_consumer_event_dedup, do:true"`
	Id           any         //
	ConsumerName any         // seller-role-assigned-consumer/...
	EventId      any         // Incoming event id
	PayloadHash  any         //
	FirstSeenAt  *gtime.Time //
}
