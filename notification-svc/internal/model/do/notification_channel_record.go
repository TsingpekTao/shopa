// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationChannelRecord is the golang structure of table notification_channel_record for DAO operations like Where/Data.
type NotificationChannelRecord struct {
	g.Meta            `orm:"table:notification_channel_record, do:true"`
	Id                any         //
	NotificationNo    any         //
	ProviderCode      any         //
	ProviderMessageId any         //
	RequestPayload    any         //
	ResponsePayload   any         //
	Status            any         //
	ErrorCode         any         //
	ErrorMessage      any         //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}
