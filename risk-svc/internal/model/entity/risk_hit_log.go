// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskHitLog is the golang structure for table risk_hit_log.
type RiskHitLog struct {
	Id        uint64      `json:"id"        orm:"id"         ` //
	HitNo     string      `json:"hitNo"     orm:"hit_no"     ` //
	UserId    uint64      `json:"userId"    orm:"user_id"    ` //
	RuleCode  string      `json:"ruleCode"  orm:"rule_code"  ` //
	BizType   string      `json:"bizType"   orm:"biz_type"   ` //
	BizNo     string      `json:"bizNo"     orm:"biz_no"     ` //
	Decision  uint        `json:"decision"  orm:"decision"   ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
}
