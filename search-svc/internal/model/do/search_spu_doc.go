// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SearchSpuDoc is the golang structure of table search_spu_doc for DAO operations like Where/Data.
type SearchSpuDoc struct {
	g.Meta            `orm:"table:search_spu_doc, do:true"`
	Id                any         //
	SpuNo             any         //
	Title             any         //
	ShopNo            any         //
	ShopName          any         //
	CategoryNo        any         //
	CoverAssetId      any         //
	CoverUrl          any         //
	MinPrice          any         //
	MaxPrice          any         //
	StockTotal        any         //
	SalesCount        any         //
	AvgScoreX100      any         //
	ReviewTotal       any         //
	ShopStatusCode    any         //
	OnShelfStatusCode any         //
	AttrsJson         any         //
	UpdatedAt         *gtime.Time //
	CreatedAt         *gtime.Time //
}
