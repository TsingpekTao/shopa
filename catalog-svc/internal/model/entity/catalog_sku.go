// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSku is the golang structure for table catalog_sku.
type CatalogSku struct {
	Id               uint64      `json:"id"               orm:"id"                  ` //
	SkuNo            string      `json:"skuNo"            orm:"sku_no"              ` //
	SpuNo            string      `json:"spuNo"            orm:"spu_no"              ` //
	ShopNo           string      `json:"shopNo"           orm:"shop_no"             ` //
	SkuName          string      `json:"skuName"          orm:"sku_name"            ` //
	SkuImageAssetId  uint64      `json:"skuImageAssetId"  orm:"sku_image_asset_id"  ` //
	SkuStatus        uint        `json:"skuStatus"        orm:"sku_status"          ` //
	StockStatus      uint        `json:"stockStatus"      orm:"stock_status"        ` //
	StockVersion     uint64      `json:"stockVersion"     orm:"stock_version"       ` // From inventory event
	LastStockEventId string      `json:"lastStockEventId" orm:"last_stock_event_id" ` //
	SalePrice        uint64      `json:"salePrice"        orm:"sale_price"          ` //
	MarketPrice      uint64      `json:"marketPrice"      orm:"market_price"        ` //
	SaleSpecsJson    string      `json:"saleSpecsJson"    orm:"sale_specs_json"     ` // Canonical sale attrs json for display
	SaleSpecsHash    string      `json:"saleSpecsHash"    orm:"sale_specs_hash"     ` // Canonical hash for unique combination
	ActiveSpecsHash  string      `json:"activeSpecsHash"  orm:"active_specs_hash"   ` //
	SortOrder        int         `json:"sortOrder"        orm:"sort_order"          ` //
	Version          uint        `json:"version"          orm:"version"             ` //
	DeletedAt        *gtime.Time `json:"deletedAt"        orm:"deleted_at"          ` //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"          ` //
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"          ` //
}
