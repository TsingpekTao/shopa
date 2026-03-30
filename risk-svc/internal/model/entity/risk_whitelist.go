// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskWhitelist is the golang structure for table risk_whitelist.
type RiskWhitelist struct {
	Id        uint64      `json:"id"        orm:"id"         ` //
	UserId    uint64      `json:"userId"    orm:"user_id"    ` //
	Reason    string      `json:"reason"    orm:"reason"     ` //
	Status    int         `json:"status"    orm:"status"     ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
