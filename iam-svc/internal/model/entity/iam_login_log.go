// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamLoginLog is the golang structure for table iam_login_log.
type IamLoginLog struct {
	Id          uint64      `json:"id"          orm:"id"          description:""`                         //
	UserId      uint64      `json:"userId"      orm:"user_id"     description:""`                         //
	Identifier  string      `json:"identifier"  orm:"identifier"  description:""`                         //
	Channel     uint        `json:"channel"     orm:"channel"     description:"1 password,2 sms,3 oauth"` // 1 password,2 sms,3 oauth
	Success     int         `json:"success"     orm:"success"     description:""`                         //
	FailReason  string      `json:"failReason"  orm:"fail_reason" description:""`                         //
	Ip          string      `json:"ip"          orm:"ip"          description:""`                         //
	Geo         string      `json:"geo"         orm:"geo"         description:""`                         //
	Ua          string      `json:"ua"          orm:"ua"          description:""`                         //
	Fingerprint string      `json:"fingerprint" orm:"fingerprint" description:""`                         //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:""`                         //
}
