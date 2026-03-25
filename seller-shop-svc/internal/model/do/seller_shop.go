// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerShop is the golang structure of table seller_shop for DAO operations like Where/Data.
type SellerShop struct {
	g.Meta              `orm:"table:seller_shop, do:true"`
	Id                  any         // Internal shop id (scope_id)
	ShopNo              any         // External business shop id
	EntityId            any         // FK seller_entity.id
	EntityNo            any         // External entity id
	ApplicationNo       any         // Source application no
	OwnerUserId         any         // IAM user id
	ShopName            any         // Requested unique shop name
	ShopNameNorm        any         // Normalized lower name
	ShopDisplayName     any         // Display name
	ShopTypeCode        any         // Shop type code
	MainCategoryIdsJson any         // Main category id list
	LogoAssetId         any         //
	BannerAssetId       any         //
	ServicePhone        any         //
	ServiceEmail        any         //
	Description         any         //
	ExtJson             any         //
	Status              any         // ShopStatus enum
	BuyerVisible        any         // 1 visible to buyers
	Version             any         // Optimistic version
	FreezeReasonCode    any         //
	FreezeReason        any         //
	CloseReasonCode     any         //
	CloseReason         any         //
	ProvisionStartedAt  *gtime.Time //
	ProvisionFinishedAt *gtime.Time //
	ActivatedAt         *gtime.Time //
	FrozenAt            *gtime.Time //
	ClosedAt            *gtime.Time //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
