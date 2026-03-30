// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskRuleVersion is the golang structure for table risk_rule_version.
type RiskRuleVersion struct {
	Id           uint64      `json:"id"           orm:"id"             ` //
	RuleCode     string      `json:"ruleCode"     orm:"rule_code"      ` //
	Version      uint        `json:"version"      orm:"version"        ` //
	RuleExprJson string      `json:"ruleExprJson" orm:"rule_expr_json" ` //
	Enabled      int         `json:"enabled"      orm:"enabled"        ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` //
}
