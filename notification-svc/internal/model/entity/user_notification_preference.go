// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserNotificationPreference is the golang structure for table user_notification_preference.
type UserNotificationPreference struct {
	Id                      uint64      `json:"id"                      orm:"id"                        ` //
	UserId                  uint64      `json:"userId"                  orm:"user_id"                   ` //
	AllowTransactionalSms   int         `json:"allowTransactionalSms"   orm:"allow_transactional_sms"   ` //
	AllowMarketingSms       int         `json:"allowMarketingSms"       orm:"allow_marketing_sms"       ` //
	AllowTransactionalEmail int         `json:"allowTransactionalEmail" orm:"allow_transactional_email" ` //
	AllowMarketingEmail     int         `json:"allowMarketingEmail"     orm:"allow_marketing_email"     ` //
	AllowPush               int         `json:"allowPush"               orm:"allow_push"                ` //
	CreatedAt               *gtime.Time `json:"createdAt"               orm:"created_at"                ` //
	UpdatedAt               *gtime.Time `json:"updatedAt"               orm:"updated_at"                ` //
}
