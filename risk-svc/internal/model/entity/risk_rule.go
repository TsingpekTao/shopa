// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskRule is the golang structure for table risk_rule.
type RiskRule struct {
	Id           uint64      `json:"id"           orm:"id"             ` //
	RuleCode     string      `json:"ruleCode"     orm:"rule_code"      ` //
	RuleName     string      `json:"ruleName"     orm:"rule_name"      ` //
	RuleExprJson string      `json:"ruleExprJson" orm:"rule_expr_json" ` //
	Priority     uint        `json:"priority"     orm:"priority"       ` //
	Enabled      int         `json:"enabled"      orm:"enabled"        ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"     ` //
}
