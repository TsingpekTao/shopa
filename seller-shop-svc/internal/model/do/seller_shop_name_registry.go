// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerShopNameRegistry is the golang structure of table seller_shop_name_registry for DAO operations like Where/Data.
type SellerShopNameRegistry struct {
	g.Meta        `orm:"table:seller_shop_name_registry, do:true"`
	Id            any         //
	ShopNameNorm  any         // Normalized unique shop name
	ShopName      any         // Original input shop name
	OwnerUserId   any         //
	ApplicationNo any         //
	Status        any         //
	ExpireAt      *gtime.Time // Reservation expiration time
	BoundShopNo   any         // Final bound shop_no when approved
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
