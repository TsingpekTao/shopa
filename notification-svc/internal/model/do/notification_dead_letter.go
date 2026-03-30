// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationDeadLetter is the golang structure of table notification_dead_letter for DAO operations like Where/Data.
type NotificationDeadLetter struct {
	g.Meta         `orm:"table:notification_dead_letter, do:true"`
	Id             any         //
	NotificationNo any         //
	ReasonCode     any         //
	PayloadJson    any         //
	CreatedAt      *gtime.Time //
}
