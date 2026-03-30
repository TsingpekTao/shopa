// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CouponUser is the golang structure of table coupon_user for DAO operations like Where/Data.
type CouponUser struct {
	g.Meta          `orm:"table:coupon_user, do:true"`
	Id              any         //
	CouponNo        any         //
	TemplateNo      any         //
	UserId          any         //
	Status          any         //
	ThresholdAmount any         //
	DiscountAmount  any         //
	ExpireAt        *gtime.Time //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
