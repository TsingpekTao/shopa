// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogCategory is the golang structure for table catalog_category.
type CatalogCategory struct {
	CategoryId   uint64      `json:"categoryId"   orm:"category_id"   ` //
	CategoryName string      `json:"categoryName" orm:"category_name" ` //
	ParentId     uint64      `json:"parentId"     orm:"parent_id"     ` //
	Level        uint        `json:"level"        orm:"level"         ` //
	Path         string      `json:"path"         orm:"path"          ` //
	SortOrder    int         `json:"sortOrder"    orm:"sort_order"    ` //
	IsLeaf       int         `json:"isLeaf"       orm:"is_leaf"       ` //
	Status       uint        `json:"status"       orm:"status"        ` // 1 ENABLED, 2 DISABLED
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
