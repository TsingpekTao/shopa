// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CartItemBackup is the golang structure for table cart_item_backup.
type CartItemBackup struct {
	Id                uint64      `json:"id"                orm:"id"                  description:""`                   //
	UserId            uint64      `json:"userId"            orm:"user_id"             description:""`                   //
	SkuNo             string      `json:"skuNo"             orm:"sku_no"              description:""`                   //
	SpuNo             string      `json:"spuNo"             orm:"spu_no"              description:""`                   //
	ShopNo            string      `json:"shopNo"            orm:"shop_no"             description:""`                   //
	Qty               uint        `json:"qty"               orm:"qty"                 description:""`                   //
	Checked           int         `json:"checked"           orm:"checked"             description:""`                   //
	Status            uint        `json:"status"            orm:"status"              description:"1=ACTIVE,2=INVALID"` // 1=ACTIVE,2=INVALID
	InvalidReasonCode string      `json:"invalidReasonCode" orm:"invalid_reason_code" description:""`                   //
	SpuTitle          string      `json:"spuTitle"          orm:"spu_title"           description:""`                   //
	SkuName           string      `json:"skuName"           orm:"sku_name"            description:""`                   //
	SkuImageAssetId   uint64      `json:"skuImageAssetId"   orm:"sku_image_asset_id"  description:""`                   //
	SalePrice         uint64      `json:"salePrice"         orm:"sale_price"          description:""`                   //
	MarketPrice       uint64      `json:"marketPrice"       orm:"market_price"        description:""`                   //
	SaleAttrsJson     string      `json:"saleAttrsJson"     orm:"sale_attrs_json"     description:""`                   //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          description:""`                   //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          description:""`                   //
}
