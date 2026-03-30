// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CouponLock is the golang structure for table coupon_lock.
type CouponLock struct {
	Id            uint64      `json:"id"            orm:"id"              ` //
	LockNo        string      `json:"lockNo"        orm:"lock_no"         ` //
	OrderNo       string      `json:"orderNo"       orm:"order_no"        ` //
	UserId        uint64      `json:"userId"        orm:"user_id"         ` //
	CouponNosJson string      `json:"couponNosJson" orm:"coupon_nos_json" ` //
	Status        uint        `json:"status"        orm:"status"          ` //
	LockExpireAt  *gtime.Time `json:"lockExpireAt"  orm:"lock_expire_at"  ` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      ` //
}
