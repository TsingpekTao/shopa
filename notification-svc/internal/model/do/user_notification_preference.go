// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserNotificationPreference is the golang structure of table user_notification_preference for DAO operations like Where/Data.
type UserNotificationPreference struct {
	g.Meta                  `orm:"table:user_notification_preference, do:true"`
	Id                      any         //
	UserId                  any         //
	AllowTransactionalSms   any         //
	AllowMarketingSms       any         //
	AllowTransactionalEmail any         //
	AllowMarketingEmail     any         //
	AllowPush               any         //
	CreatedAt               *gtime.Time //
	UpdatedAt               *gtime.Time //
}
