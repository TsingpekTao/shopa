// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskRuleVersion is the golang structure of table risk_rule_version for DAO operations like Where/Data.
type RiskRuleVersion struct {
	g.Meta       `orm:"table:risk_rule_version, do:true"`
	Id           any         //
	RuleCode     any         //
	Version      any         //
	RuleExprJson any         //
	Enabled      any         //
	CreatedAt    *gtime.Time //
}
