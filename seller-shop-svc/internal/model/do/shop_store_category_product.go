// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShopStoreCategoryProduct is the golang structure of table shop_store_category_product for DAO operations like Where/Data.
type ShopStoreCategoryProduct struct {
	g.Meta          `orm:"table:shop_store_category_product, do:true"`
	Id              any         //
	ShopNo          any         //
	SpuNo           any         //
	StoreCategoryId any         //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
