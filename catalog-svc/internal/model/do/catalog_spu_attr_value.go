// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSpuAttrValue is the golang structure of table catalog_spu_attr_value for DAO operations like Where/Data.
type CatalogSpuAttrValue struct {
	g.Meta    `orm:"table:catalog_spu_attr_value, do:true"`
	Id        any         //
	SpuNo     any         //
	AttrCode  any         //
	AttrName  any         //
	AttrScope any         //
	AttrValue any         //
	SortOrder any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
