// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationChannelRecord is the golang structure for table notification_channel_record.
type NotificationChannelRecord struct {
	Id                uint64      `json:"id"                orm:"id"                  ` //
	NotificationNo    string      `json:"notificationNo"    orm:"notification_no"     ` //
	ProviderCode      string      `json:"providerCode"      orm:"provider_code"       ` //
	ProviderMessageId string      `json:"providerMessageId" orm:"provider_message_id" ` //
	RequestPayload    string      `json:"requestPayload"    orm:"request_payload"     ` //
	ResponsePayload   string      `json:"responsePayload"   orm:"response_payload"    ` //
	Status            uint        `json:"status"            orm:"status"              ` //
	ErrorCode         string      `json:"errorCode"         orm:"error_code"          ` //
	ErrorMessage      string      `json:"errorMessage"      orm:"error_message"       ` //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          ` //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          ` //
}
