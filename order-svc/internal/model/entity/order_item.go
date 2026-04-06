// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItem is the golang structure for table order_item.
type OrderItem struct {
	Id              uint64      `json:"id"              orm:"id"                 ` //
	ItemNo          string      `json:"itemNo"          orm:"item_no"            ` //
	OrderNo         string      `json:"orderNo"         orm:"order_no"           ` //
	SubOrderNo      string      `json:"subOrderNo"      orm:"sub_order_no"       ` //
	ShopNo          string      `json:"shopNo"          orm:"shop_no"            ` //
	SpuNo           string      `json:"spuNo"           orm:"spu_no"             ` //
	SkuNo           string      `json:"skuNo"           orm:"sku_no"             ` //
	SpuTitle        string      `json:"spuTitle"        orm:"spu_title"          ` //
	SkuName         string      `json:"skuName"         orm:"sku_name"           ` //
	SkuImageAssetId uint64      `json:"skuImageAssetId" orm:"sku_image_asset_id" ` //
	Qty             uint        `json:"qty"             orm:"qty"                ` //
	SalePrice       uint64      `json:"salePrice"       orm:"sale_price"         ` //
	MarketPrice     uint64      `json:"marketPrice"     orm:"market_price"       ` //
	SaleAttrsJson   string      `json:"saleAttrsJson"   orm:"sale_attrs_json"    ` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"         ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"         ` //
}
