// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamSmsLog is the golang structure for table iam_sms_log.
type IamSmsLog struct {
	Id          uint64      `json:"id"          orm:"id"          description:""`                                          //
	Scene       uint        `json:"scene"       orm:"scene"       description:"1 register,2 login,3 mfa,4 reset_password"` // 1 register,2 login,3 mfa,4 reset_password
	Target      string      `json:"target"      orm:"target"      description:""`                                          //
	Provider    string      `json:"provider"    orm:"provider"    description:""`                                          //
	BizId       string      `json:"bizId"       orm:"biz_id"      description:""`                                          //
	Ip          string      `json:"ip"          orm:"ip"          description:""`                                          //
	Ua          string      `json:"ua"          orm:"ua"          description:""`                                          //
	Fingerprint string      `json:"fingerprint" orm:"fingerprint" description:""`                                          //
	Success     int         `json:"success"     orm:"success"     description:""`                                          //
	ErrorCode   int         `json:"errorCode"   orm:"error_code"  description:""`                                          //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:""`                                          //
}
