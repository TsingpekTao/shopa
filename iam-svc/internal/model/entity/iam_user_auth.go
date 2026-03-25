// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserAuth is the golang structure for table iam_user_auth.
type IamUserAuth struct {
	UserId               uint64      `json:"userId"               orm:"user_id"                description:"Snowflake user id"`                           // Snowflake user id
	Phone                string      `json:"phone"                orm:"phone"                  description:"Phone number"`                                // Phone number
	Email                string      `json:"email"                orm:"email"                  description:"Email address"`                               // Email address
	PasswordHash         string      `json:"passwordHash"         orm:"password_hash"          description:"Password hash"`                               // Password hash
	PasswordSalt         string      `json:"passwordSalt"         orm:"password_salt"          description:"Password salt"`                               // Password salt
	PasswordAlgo         string      `json:"passwordAlgo"         orm:"password_algo"          description:"Password hash algorithm"`                     // Password hash algorithm
	PasswordVer          uint        `json:"passwordVer"          orm:"password_ver"           description:"Password algorithm version"`                  // Password algorithm version
	AccountStatus        uint        `json:"accountStatus"        orm:"account_status"         description:"1 active,2 locked,3 disabled,4 review"`       // 1 active,2 locked,3 disabled,4 review
	FailedLoginCount     uint        `json:"failedLoginCount"     orm:"failed_login_count"     description:"Fail count for audit"`                        // Fail count for audit
	LockedUntil          *gtime.Time `json:"lockedUntil"          orm:"locked_until"           description:"Account lock expiry time"`                    // Account lock expiry time
	TokenVersion         uint        `json:"tokenVersion"         orm:"token_version"          description:"Token version for global token invalidation"` // Token version for global token invalidation
	LastLoginAt          *gtime.Time `json:"lastLoginAt"          orm:"last_login_at"          description:"Last login time"`                             // Last login time
	LastLoginIp          string      `json:"lastLoginIp"          orm:"last_login_ip"          description:"Last login ip"`                               // Last login ip
	LastLoginGeo         string      `json:"lastLoginGeo"         orm:"last_login_geo"         description:"Last login geo"`                              // Last login geo
	LastLoginUa          string      `json:"lastLoginUa"          orm:"last_login_ua"          description:"Last login user-agent"`                       // Last login user-agent
	LastLoginFingerprint string      `json:"lastLoginFingerprint" orm:"last_login_fingerprint" description:"Last login fingerprint"`                      // Last login fingerprint
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             description:""`                                            //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:""`                                            //
}
