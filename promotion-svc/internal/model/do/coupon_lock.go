// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CouponLock is the golang structure of table coupon_lock for DAO operations like Where/Data.
type CouponLock struct {
	g.Meta        `orm:"table:coupon_lock, do:true"`
	Id            any         //
	LockNo        any         //
	OrderNo       any         //
	UserId        any         //
	CouponNosJson any         //
	Status        any         //
	LockExpireAt  *gtime.Time //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
