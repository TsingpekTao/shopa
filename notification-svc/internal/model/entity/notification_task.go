// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationTask is the golang structure for table notification_task.
type NotificationTask struct {
	Id                  uint64      `json:"id"                  orm:"id"                    ` //
	NotificationNo      string      `json:"notificationNo"      orm:"notification_no"       ` //
	UserId              uint64      `json:"userId"              orm:"user_id"               ` //
	TargetAddress       string      `json:"targetAddress"       orm:"target_address"        ` //
	TemplateCode        string      `json:"templateCode"        orm:"template_code"         ` //
	Channel             uint        `json:"channel"             orm:"channel"               ` //
	BizType             uint        `json:"bizType"             orm:"biz_type"              ` //
	TemplateParamsJson  string      `json:"templateParamsJson"  orm:"template_params_json"  ` //
	Status              uint        `json:"status"              orm:"status"                ` //
	BlockedByPreference int         `json:"blockedByPreference" orm:"blocked_by_preference" ` //
	BlockedByFrequency  int         `json:"blockedByFrequency"  orm:"blocked_by_frequency"  ` //
	RetryCount          uint        `json:"retryCount"          orm:"retry_count"           ` //
	NextRetryAt         *gtime.Time `json:"nextRetryAt"         orm:"next_retry_at"         ` //
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            ` //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            ` //
}
