// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationTask is the golang structure of table notification_task for DAO operations like Where/Data.
type NotificationTask struct {
	g.Meta              `orm:"table:notification_task, do:true"`
	Id                  any         //
	NotificationNo      any         //
	UserId              any         //
	TargetAddress       any         //
	TemplateCode        any         //
	Channel             any         //
	BizType             any         //
	TemplateParamsJson  any         //
	Status              any         //
	BlockedByPreference any         //
	BlockedByFrequency  any         //
	RetryCount          any         //
	NextRetryAt         *gtime.Time //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
