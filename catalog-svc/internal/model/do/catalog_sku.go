// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSku is the golang structure of table catalog_sku for DAO operations like Where/Data.
type CatalogSku struct {
	g.Meta           `orm:"table:catalog_sku, do:true"`
	Id               any         //
	SkuNo            any         //
	SpuNo            any         //
	ShopNo           any         //
	SkuName          any         //
	SkuImageAssetId  any         //
	SkuStatus        any         //
	StockStatus      any         //
	StockVersion     any         // From inventory event
	LastStockEventId any         //
	SalePrice        any         //
	MarketPrice      any         //
	SaleSpecsJson    any         // Canonical sale attrs json for display
	SaleSpecsHash    any         // Canonical hash for unique combination
	ActiveSpecsHash  any         //
	SortOrder        any         //
	Version          any         //
	DeletedAt        *gtime.Time //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
}
