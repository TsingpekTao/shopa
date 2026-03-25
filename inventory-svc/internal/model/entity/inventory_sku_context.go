// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventorySkuContext 是 inventory_sku_context 表的 Go 结构体。
type InventorySkuContext struct {
	Id        uint64      `json:"id"        orm:"id"         ` //
	SkuNo     string      `json:"skuNo"     orm:"sku_no"     ` //
	SpuNo     string      `json:"spuNo"     orm:"spu_no"     ` //
	ShopNo    string      `json:"shopNo"    orm:"shop_no"    ` //
	Enabled   int         `json:"enabled"   orm:"enabled"    ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
