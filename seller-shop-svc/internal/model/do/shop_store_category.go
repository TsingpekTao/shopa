// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShopStoreCategory is the golang structure of table shop_store_category for DAO operations like Where/Data.
type ShopStoreCategory struct {
	g.Meta       `orm:"table:shop_store_category, do:true"`
	Id           any         //
	ShopNo       any         //
	ParentId     any         //
	Name         any         //
	Level        any         //
	SortOrder    any         //
	IsVisible    any         //
	IsDeleted    any         //
	ProductCount any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
