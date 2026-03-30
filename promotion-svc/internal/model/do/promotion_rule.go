// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionRule is the golang structure of table promotion_rule for DAO operations like Where/Data.
type PromotionRule struct {
	g.Meta     `orm:"table:promotion_rule, do:true"`
	Id         any         //
	RuleCode   any         //
	CampaignNo any         //
	RuleJson   any         //
	Priority   any         //
	Status     any         //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
