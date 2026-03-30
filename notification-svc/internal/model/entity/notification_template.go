// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationTemplate is the golang structure for table notification_template.
type NotificationTemplate struct {
	Id              uint64      `json:"id"              orm:"id"               ` //
	TemplateCode    string      `json:"templateCode"    orm:"template_code"    ` //
	Channel         uint        `json:"channel"         orm:"channel"          ` //
	BizType         uint        `json:"bizType"         orm:"biz_type"         ` //
	TitleTemplate   string      `json:"titleTemplate"   orm:"title_template"   ` //
	ContentTemplate string      `json:"contentTemplate" orm:"content_template" ` //
	Status          int         `json:"status"          orm:"status"           ` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` //
}
