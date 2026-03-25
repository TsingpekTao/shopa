// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerShopNameRegistry is the golang structure for table seller_shop_name_registry.
type SellerShopNameRegistry struct {
	Id            uint64      `json:"id"            orm:"id"             ` //
	ShopNameNorm  string      `json:"shopNameNorm"  orm:"shop_name_norm" ` // Normalized unique shop name
	ShopName      string      `json:"shopName"      orm:"shop_name"      ` // Original input shop name
	OwnerUserId   uint64      `json:"ownerUserId"   orm:"owner_user_id"  ` //
	ApplicationNo string      `json:"applicationNo" orm:"application_no" ` //
	Status        uint        `json:"status"        orm:"status"         ` //
	ExpireAt      *gtime.Time `json:"expireAt"      orm:"expire_at"      ` // Reservation expiration time
	BoundShopNo   string      `json:"boundShopNo"   orm:"bound_shop_no"  ` // Final bound shop_no when approved
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` //
}
