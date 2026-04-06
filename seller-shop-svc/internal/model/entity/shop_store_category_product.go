// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ShopStoreCategoryProduct is the golang structure for table shop_store_category_product.
type ShopStoreCategoryProduct struct {
	Id              uint64      `json:"id"              orm:"id"                ` //
	ShopNo          string      `json:"shopNo"          orm:"shop_no"           ` //
	SpuNo           string      `json:"spuNo"           orm:"spu_no"            ` //
	StoreCategoryId uint64      `json:"storeCategoryId" orm:"store_category_id" ` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"        ` //
}
