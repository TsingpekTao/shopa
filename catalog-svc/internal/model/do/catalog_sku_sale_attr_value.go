// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSkuSaleAttrValue is the golang structure of table catalog_sku_sale_attr_value for DAO operations like Where/Data.
type CatalogSkuSaleAttrValue struct {
	g.Meta    `orm:"table:catalog_sku_sale_attr_value, do:true"`
	Id        any         //
	SkuNo     any         //
	SpuNo     any         //
	AttrCode  any         //
	AttrName  any         //
	AttrValue any         //
	SortOrder any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
