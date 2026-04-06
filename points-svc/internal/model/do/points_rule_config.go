// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsRuleConfig is the golang structure of table points_rule_config for DAO operations like Where/Data.
type PointsRuleConfig struct {
	g.Meta              `orm:"table:points_rule_config, do:true"`
	RuleCode            any         //
	RuleName            any         //
	StatusCode          any         //
	MinOrderAmountCent  any         //
	MaxDeductionRateBps any         // Max deduction rate in basis points
	DeductPointsPerCent any         // How many points are required for one cent discount
	GrantPointsPerCent  any         // How many points are granted for one cent paid amount
	RefundGraceDays     any         //
	EffectiveAt         *gtime.Time //
	ExpireAt            *gtime.Time //
	RuleSnapshotJson    any         //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
