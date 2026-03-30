// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionCampaign is the golang structure of table promotion_campaign for DAO operations like Where/Data.
type PromotionCampaign struct {
	g.Meta     `orm:"table:promotion_campaign, do:true"`
	Id         any         //
	CampaignNo any         //
	Name       any         //
	RuleJson   any         //
	StartAt    *gtime.Time //
	EndAt      *gtime.Time //
	Status     any         //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
