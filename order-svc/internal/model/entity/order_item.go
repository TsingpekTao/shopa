// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItem is the golang structure for table order_item.
type OrderItem struct {
	Id              uint64      `json:"id"              orm:"id"                 description:""` //
	ItemNo          string      `json:"itemNo"          orm:"item_no"            description:""` //
	OrderNo         string      `json:"orderNo"         orm:"order_no"           description:""` //
	SubOrderNo      string      `json:"subOrderNo"      orm:"sub_order_no"       description:""` //
	ShopNo          string      `json:"shopNo"          orm:"shop_no"            description:""` //
	SpuNo           string      `json:"spuNo"           orm:"spu_no"             description:""` //
	SkuNo           string      `json:"skuNo"           orm:"sku_no"             description:""` //
	SpuTitle        string      `json:"spuTitle"        orm:"spu_title"          description:""` //
	SkuName         string      `json:"skuName"         orm:"sku_name"           description:""` //
	SkuImageAssetId uint64      `json:"skuImageAssetId" orm:"sku_image_asset_id" description:""` //
	Qty             uint        `json:"qty"             orm:"qty"                description:""` //
	SalePrice       uint64      `json:"salePrice"       orm:"sale_price"         description:""` //
	MarketPrice     uint64      `json:"marketPrice"     orm:"market_price"       description:""` //
	SaleAttrsJson   string      `json:"saleAttrsJson"   orm:"sale_attrs_json"    description:""` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"         description:""` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"         description:""` //
}
