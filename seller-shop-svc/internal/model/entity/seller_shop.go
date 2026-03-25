// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerShop is the golang structure for table seller_shop.
type SellerShop struct {
	Id                  uint64      `json:"id"                  orm:"id"                     ` // Internal shop id (scope_id)
	ShopNo              string      `json:"shopNo"              orm:"shop_no"                ` // External business shop id
	EntityId            uint64      `json:"entityId"            orm:"entity_id"              ` // FK seller_entity.id
	EntityNo            string      `json:"entityNo"            orm:"entity_no"              ` // External entity id
	ApplicationNo       string      `json:"applicationNo"       orm:"application_no"         ` // Source application no
	OwnerUserId         uint64      `json:"ownerUserId"         orm:"owner_user_id"          ` // IAM user id
	ShopName            string      `json:"shopName"            orm:"shop_name"              ` // Requested unique shop name
	ShopNameNorm        string      `json:"shopNameNorm"        orm:"shop_name_norm"         ` // Normalized lower name
	ShopDisplayName     string      `json:"shopDisplayName"     orm:"shop_display_name"      ` // Display name
	ShopTypeCode        string      `json:"shopTypeCode"        orm:"shop_type_code"         ` // Shop type code
	MainCategoryIdsJson string      `json:"mainCategoryIdsJson" orm:"main_category_ids_json" ` // Main category id list
	LogoAssetId         uint64      `json:"logoAssetId"         orm:"logo_asset_id"          ` //
	BannerAssetId       uint64      `json:"bannerAssetId"       orm:"banner_asset_id"        ` //
	ServicePhone        string      `json:"servicePhone"        orm:"service_phone"          ` //
	ServiceEmail        string      `json:"serviceEmail"        orm:"service_email"          ` //
	Description         string      `json:"description"         orm:"description"            ` //
	ExtJson             string      `json:"extJson"             orm:"ext_json"               ` //
	Status              uint        `json:"status"              orm:"status"                 ` // ShopStatus enum
	BuyerVisible        int         `json:"buyerVisible"        orm:"buyer_visible"          ` // 1 visible to buyers
	Version             uint        `json:"version"             orm:"version"                ` // Optimistic version
	FreezeReasonCode    string      `json:"freezeReasonCode"    orm:"freeze_reason_code"     ` //
	FreezeReason        string      `json:"freezeReason"        orm:"freeze_reason"          ` //
	CloseReasonCode     string      `json:"closeReasonCode"     orm:"close_reason_code"      ` //
	CloseReason         string      `json:"closeReason"         orm:"close_reason"           ` //
	ProvisionStartedAt  *gtime.Time `json:"provisionStartedAt"  orm:"provision_started_at"   ` //
	ProvisionFinishedAt *gtime.Time `json:"provisionFinishedAt" orm:"provision_finished_at"  ` //
	ActivatedAt         *gtime.Time `json:"activatedAt"         orm:"activated_at"           ` //
	FrozenAt            *gtime.Time `json:"frozenAt"            orm:"frozen_at"              ` //
	ClosedAt            *gtime.Time `json:"closedAt"            orm:"closed_at"              ` //
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"             ` //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"             ` //
}
