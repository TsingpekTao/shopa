// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogAttributeTemplate is the golang structure of table catalog_attribute_template for DAO operations like Where/Data.
type CatalogAttributeTemplate struct {
	g.Meta         `orm:"table:catalog_attribute_template, do:true"`
	Id             any         //
	CategoryId     any         //
	AttrCode       any         //
	AttrName       any         //
	AttrScope      any         // 1 SPU, 2 SKU_SALE, 3 EXT
	ValueType      any         // TEXT/NUMBER/BOOL/ENUM
	RequiredFlag   any         //
	SearchableFlag any         //
	FilterableFlag any         //
	OptionsJson    any         // Enum options if needed
	SortOrder      any         //
	Status         any         // 1 ENABLED, 2 DISABLED
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
