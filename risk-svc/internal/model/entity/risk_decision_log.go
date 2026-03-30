// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskDecisionLog is the golang structure for table risk_decision_log.
type RiskDecisionLog struct {
	Id               uint64      `json:"id"               orm:"id"                 ` //
	DecisionNo       string      `json:"decisionNo"       orm:"decision_no"        ` //
	UserId           uint64      `json:"userId"           orm:"user_id"            ` //
	BizType          string      `json:"bizType"          orm:"biz_type"           ` //
	BizNo            string      `json:"bizNo"            orm:"biz_no"             ` //
	Decision         uint        `json:"decision"         orm:"decision"           ` //
	RiskScore        uint        `json:"riskScore"        orm:"risk_score"         ` //
	MatchedRuleCodes string      `json:"matchedRuleCodes" orm:"matched_rule_codes" ` //
	ReasonCodes      string      `json:"reasonCodes"      orm:"reason_codes"       ` //
	Degraded         int         `json:"degraded"         orm:"degraded"           ` //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         ` //
}
