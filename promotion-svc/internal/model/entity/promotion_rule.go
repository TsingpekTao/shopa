// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionRule is the golang structure for table promotion_rule.
type PromotionRule struct {
	Id         uint64      `json:"id"         orm:"id"          ` //
	RuleCode   string      `json:"ruleCode"   orm:"rule_code"   ` //
	CampaignNo string      `json:"campaignNo" orm:"campaign_no" ` //
	RuleJson   string      `json:"ruleJson"   orm:"rule_json"   ` //
	Priority   uint        `json:"priority"   orm:"priority"    ` //
	Status     uint        `json:"status"     orm:"status"      ` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` //
}
