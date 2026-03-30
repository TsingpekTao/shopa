// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskDecisionLog is the golang structure of table risk_decision_log for DAO operations like Where/Data.
type RiskDecisionLog struct {
	g.Meta           `orm:"table:risk_decision_log, do:true"`
	Id               any         //
	DecisionNo       any         //
	UserId           any         //
	BizType          any         //
	BizNo            any         //
	Decision         any         //
	RiskScore        any         //
	MatchedRuleCodes any         //
	ReasonCodes      any         //
	Degraded         any         //
	CreatedAt        *gtime.Time //
}
