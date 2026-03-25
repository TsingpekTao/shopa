// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamLoginLog is the golang structure of table iam_login_log for DAO operations like Where/Data.
type IamLoginLog struct {
	g.Meta      `orm:"table:iam_login_log, do:true"`
	Id          any         //
	UserId      any         //
	Identifier  any         //
	Channel     any         // 1 password,2 sms,3 oauth
	Success     any         //
	FailReason  any         //
	Ip          any         //
	Geo         any         //
	Ua          any         //
	Fingerprint any         //
	CreatedAt   *gtime.Time //
}
