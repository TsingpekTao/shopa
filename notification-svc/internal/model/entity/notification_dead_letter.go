// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotificationDeadLetter is the golang structure for table notification_dead_letter.
type NotificationDeadLetter struct {
	Id             uint64      `json:"id"             orm:"id"              ` //
	NotificationNo string      `json:"notificationNo" orm:"notification_no" ` //
	ReasonCode     string      `json:"reasonCode"     orm:"reason_code"     ` //
	PayloadJson    string      `json:"payloadJson"    orm:"payload_json"    ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
}
