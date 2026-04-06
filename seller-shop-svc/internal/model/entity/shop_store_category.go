// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ShopStoreCategory is the golang structure for table shop_store_category.
type ShopStoreCategory struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	ShopNo       string      `json:"shopNo"       orm:"shop_no"       ` //
	ParentId     uint64      `json:"parentId"     orm:"parent_id"     ` //
	Name         string      `json:"name"         orm:"name"          ` //
	Level        uint        `json:"level"        orm:"level"         ` //
	SortOrder    int         `json:"sortOrder"    orm:"sort_order"    ` //
	IsVisible    int         `json:"isVisible"    orm:"is_visible"    ` //
	IsDeleted    int         `json:"isDeleted"    orm:"is_deleted"    ` //
	ProductCount uint        `json:"productCount" orm:"product_count" ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
