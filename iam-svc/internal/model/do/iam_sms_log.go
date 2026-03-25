// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamSmsLog is the golang structure of table iam_sms_log for DAO operations like Where/Data.
type IamSmsLog struct {
	g.Meta      `orm:"table:iam_sms_log, do:true"`
	Id          any         //
	Scene       any         // 1 register,2 login,3 mfa,4 reset_password
	Target      any         //
	Provider    any         //
	BizId       any         //
	Ip          any         //
	Ua          any         //
	Fingerprint any         //
	Success     any         //
	ErrorCode   any         //
	CreatedAt   *gtime.Time //
}
