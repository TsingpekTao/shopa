// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsRuleConfig is the golang structure for table points_rule_config.
type PointsRuleConfig struct {
	RuleCode            string      `json:"ruleCode"            orm:"rule_code"              ` //
	RuleName            string      `json:"ruleName"            orm:"rule_name"              ` //
	StatusCode          string      `json:"statusCode"          orm:"status_code"            ` //
	MinOrderAmountCent  int64       `json:"minOrderAmountCent"  orm:"min_order_amount_cent"  ` //
	MaxDeductionRateBps uint        `json:"maxDeductionRateBps" orm:"max_deduction_rate_bps" ` // Max deduction rate in basis points
	DeductPointsPerCent uint64      `json:"deductPointsPerCent" orm:"deduct_points_per_cent" ` // How many points are required for one cent discount
	GrantPointsPerCent  uint64      `json:"grantPointsPerCent"  orm:"grant_points_per_cent"  ` // How many points are granted for one cent paid amount
	RefundGraceDays     uint        `json:"refundGraceDays"     orm:"refund_grace_days"      ` //
	EffectiveAt         *gtime.Time `json:"effectiveAt"         orm:"effective_at"           ` //
	ExpireAt            *gtime.Time `json:"expireAt"            orm:"expire_at"              ` //
	RuleSnapshotJson    string      `json:"ruleSnapshotJson"    orm:"rule_snapshot_json"     ` //
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"             ` //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"             ` //
}
