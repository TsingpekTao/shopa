// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaConsumerEventDedup is the golang structure of table media_consumer_event_dedup for DAO operations like Where/Data.
type MediaConsumerEventDedup struct {
	g.Meta       `orm:"table:media_consumer_event_dedup, do:true"`
	Id           any         //
	ConsumerName any         //
	EventId      any         //
	PayloadHash  any         //
	FirstSeenAt  *gtime.Time //
}
