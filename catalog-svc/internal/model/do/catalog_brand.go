// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogBrand is the golang structure of table catalog_brand for DAO operations like Where/Data.
type CatalogBrand struct {
	g.Meta    `orm:"table:catalog_brand, do:true"`
	Id        any         //
	BrandNo   any         //
	BrandName any         //
	Status    any         // 1 ENABLED, 2 DISABLED
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
