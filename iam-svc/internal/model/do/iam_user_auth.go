// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserAuth is the golang structure of table iam_user_auth for DAO operations like Where/Data.
type IamUserAuth struct {
	g.Meta               `orm:"table:iam_user_auth, do:true"`
	UserId               any         // Snowflake user id
	Phone                any         // Phone number
	Email                any         // Email address
	PasswordHash         any         // Password hash
	PasswordSalt         any         // Password salt
	PasswordAlgo         any         // Password hash algorithm
	PasswordVer          any         // Password algorithm version
	AccountStatus        any         // 1 active,2 locked,3 disabled,4 review
	FailedLoginCount     any         // Fail count for audit
	LockedUntil          *gtime.Time // Account lock expiry time
	TokenVersion         any         // Token version for global token invalidation
	LastLoginAt          *gtime.Time // Last login time
	LastLoginIp          any         // Last login ip
	LastLoginGeo         any         // Last login geo
	LastLoginUa          any         // Last login user-agent
	LastLoginFingerprint any         // Last login fingerprint
	CreatedAt            *gtime.Time //
	UpdatedAt            *gtime.Time //
}
