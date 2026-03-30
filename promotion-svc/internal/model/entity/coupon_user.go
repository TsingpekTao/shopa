// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CouponUser is the golang structure for table coupon_user.
type CouponUser struct {
	Id              uint64      `json:"id"              orm:"id"               ` //
	CouponNo        string      `json:"couponNo"        orm:"coupon_no"        ` //
	TemplateNo      string      `json:"templateNo"      orm:"template_no"      ` //
	UserId          uint64      `json:"userId"          orm:"user_id"          ` //
	Status          uint        `json:"status"          orm:"status"           ` //
	ThresholdAmount uint64      `json:"thresholdAmount" orm:"threshold_amount" ` //
	DiscountAmount  uint64      `json:"discountAmount"  orm:"discount_amount"  ` //
	ExpireAt        *gtime.Time `json:"expireAt"        orm:"expire_at"        ` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` //
}
