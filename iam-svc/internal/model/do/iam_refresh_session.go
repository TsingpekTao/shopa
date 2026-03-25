// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamRefreshSession is the golang structure of table iam_refresh_session for DAO operations like Where/Data.
type IamRefreshSession struct {
	g.Meta           `orm:"table:iam_refresh_session, do:true"`
	Sid              any         // Session id
	UserId           any         //
	RefreshTokenHash any         //
	UaHash           any         //
	Ip               any         //
	Geo              any         //
	Fingerprint      any         //
	CreatedAt        *gtime.Time //
	ExpiresAt        *gtime.Time //
	RevokedAt        *gtime.Time //
	ReplacedBySid    any         //
}
