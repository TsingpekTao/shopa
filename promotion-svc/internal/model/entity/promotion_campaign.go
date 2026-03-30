// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionCampaign is the golang structure for table promotion_campaign.
type PromotionCampaign struct {
	Id         uint64      `json:"id"         orm:"id"          ` //
	CampaignNo string      `json:"campaignNo" orm:"campaign_no" ` //
	Name       string      `json:"name"       orm:"name"        ` //
	RuleJson   string      `json:"ruleJson"   orm:"rule_json"   ` //
	StartAt    *gtime.Time `json:"startAt"    orm:"start_at"    ` //
	EndAt      *gtime.Time `json:"endAt"      orm:"end_at"      ` //
	Status     uint        `json:"status"     orm:"status"      ` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` //
}
