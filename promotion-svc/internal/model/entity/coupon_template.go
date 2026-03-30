// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CouponTemplate is the golang structure for table coupon_template.
type CouponTemplate struct {
	Id              uint64      `json:"id"              orm:"id"               ` //
	TemplateNo      string      `json:"templateNo"      orm:"template_no"      ` //
	CampaignNo      string      `json:"campaignNo"      orm:"campaign_no"      ` //
	ThresholdAmount uint64      `json:"thresholdAmount" orm:"threshold_amount" ` //
	DiscountAmount  uint64      `json:"discountAmount"  orm:"discount_amount"  ` //
	TotalCount      uint        `json:"totalCount"      orm:"total_count"      ` //
	IssuedCount     uint        `json:"issuedCount"     orm:"issued_count"     ` //
	Status          uint        `json:"status"          orm:"status"           ` //
	ExpireAt        *gtime.Time `json:"expireAt"        orm:"expire_at"        ` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` //
}
