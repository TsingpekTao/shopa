// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogCategory is the golang structure of table catalog_category for DAO operations like Where/Data.
type CatalogCategory struct {
	g.Meta       `orm:"table:catalog_category, do:true"`
	CategoryId   any         //
	CategoryName any         //
	ParentId     any         //
	Level        any         //
	Path         any         //
	SortOrder    any         //
	IsLeaf       any         //
	Status       any         // 1 ENABLED, 2 DISABLED
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
