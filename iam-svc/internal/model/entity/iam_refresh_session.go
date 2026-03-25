// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamRefreshSession is the golang structure for table iam_refresh_session.
type IamRefreshSession struct {
	Sid              string      `json:"sid"              orm:"sid"                description:"Session id"` // Session id
	UserId           uint64      `json:"userId"           orm:"user_id"            description:""`           //
	RefreshTokenHash string      `json:"refreshTokenHash" orm:"refresh_token_hash" description:""`           //
	UaHash           string      `json:"uaHash"           orm:"ua_hash"            description:""`           //
	Ip               string      `json:"ip"               orm:"ip"                 description:""`           //
	Geo              string      `json:"geo"              orm:"geo"                description:""`           //
	Fingerprint      string      `json:"fingerprint"      orm:"fingerprint"        description:""`           //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:""`           //
	ExpiresAt        *gtime.Time `json:"expiresAt"        orm:"expires_at"         description:""`           //
	RevokedAt        *gtime.Time `json:"revokedAt"        orm:"revoked_at"         description:""`           //
	ReplacedBySid    string      `json:"replacedBySid"    orm:"replaced_by_sid"    description:""`           //
}
