// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskRule is the golang structure of table risk_rule for DAO operations like Where/Data.
type RiskRule struct {
	g.Meta       `orm:"table:risk_rule, do:true"`
	Id           any         //
	RuleCode     any         //
	RuleName     any         //
	RuleExprJson any         //
	Priority     any         //
	Enabled      any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
