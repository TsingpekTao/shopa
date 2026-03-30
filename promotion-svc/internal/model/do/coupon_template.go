// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CouponTemplate is the golang structure of table coupon_template for DAO operations like Where/Data.
type CouponTemplate struct {
	g.Meta          `orm:"table:coupon_template, do:true"`
	Id              any         //
	TemplateNo      any         //
	CampaignNo      any         //
	ThresholdAmount any         //
	DiscountAmount  any         //
	TotalCount      any         //
	IssuedCount     any         //
	Status          any         //
	ExpireAt        *gtime.Time //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
