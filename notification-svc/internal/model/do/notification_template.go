// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationTemplate is the golang structure of table notification_template for DAO operations like Where/Data.
type NotificationTemplate struct {
	g.Meta          `orm:"table:notification_template, do:true"`
	Id              any         //
	TemplateCode    any         //
	Channel         any         //
	BizType         any         //
	TitleTemplate   any         //
	ContentTemplate any         //
	Status          any         //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
