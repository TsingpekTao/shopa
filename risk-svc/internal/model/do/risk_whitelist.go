// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskWhitelist is the golang structure of table risk_whitelist for DAO operations like Where/Data.
type RiskWhitelist struct {
	g.Meta    `orm:"table:risk_whitelist, do:true"`
	Id        any         //
	UserId    any         //
	Reason    any         //
	Status    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
