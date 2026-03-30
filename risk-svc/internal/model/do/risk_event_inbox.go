// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskEventInbox is the golang structure of table risk_event_inbox for DAO operations like Where/Data.
type RiskEventInbox struct {
	g.Meta      `orm:"table:risk_event_inbox, do:true"`
	Id          any         //
	EventId     any         //
	EventType   any         //
	UserId      any         //
	BizNo       any         //
	PayloadJson any         //
	OccurredAt  *gtime.Time //
	CreatedAt   *gtime.Time //
}
