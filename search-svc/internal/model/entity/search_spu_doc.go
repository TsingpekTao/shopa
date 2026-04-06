// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SearchSpuDoc is the golang structure for table search_spu_doc.
type SearchSpuDoc struct {
	Id                    uint64      `json:"id"                    orm:"id"                       ` //
	SpuNo                 string      `json:"spuNo"                 orm:"spu_no"                   ` //
	Title                 string      `json:"title"                 orm:"title"                    ` //
	ShopNo                string      `json:"shopNo"                orm:"shop_no"                  ` //
	ShopName              string      `json:"shopName"              orm:"shop_name"                ` //
	CategoryNo            string      `json:"categoryNo"            orm:"category_no"              ` //
	StoreCategoryId       uint64      `json:"storeCategoryId"       orm:"store_category_id"        ` //
	StoreCategoryL1       uint64      `json:"storeCategoryL1"       orm:"store_category_l1"        ` //
	StoreCategoryL2       uint64      `json:"storeCategoryL2"       orm:"store_category_l2"        ` //
	StoreCategoryPathJson string      `json:"storeCategoryPathJson" orm:"store_category_path_json" ` //
	CoverAssetId          uint64      `json:"coverAssetId"          orm:"cover_asset_id"           ` //
	CoverUrl              string      `json:"coverUrl"              orm:"cover_url"                ` //
	MinPrice              uint64      `json:"minPrice"              orm:"min_price"                ` //
	MaxPrice              uint64      `json:"maxPrice"              orm:"max_price"                ` //
	StockTotal            uint64      `json:"stockTotal"            orm:"stock_total"              ` //
	SalesCount            uint64      `json:"salesCount"            orm:"sales_count"              ` //
	AvgScoreX100          uint64      `json:"avgScoreX100"          orm:"avg_score_x100"           ` //
	ReviewTotal           uint64      `json:"reviewTotal"           orm:"review_total"             ` //
	ShopStatusCode        string      `json:"shopStatusCode"        orm:"shop_status_code"         ` //
	OnShelfStatusCode     string      `json:"onShelfStatusCode"     orm:"on_shelf_status_code"     ` //
	AttrsJson             string      `json:"attrsJson"             orm:"attrs_json"               ` //
	SourceVersion         uint64      `json:"sourceVersion"         orm:"source_version"           ` //
	SourceUpdatedAt       *gtime.Time `json:"sourceUpdatedAt"       orm:"source_updated_at"        ` //
	Deleted               uint        `json:"deleted"               orm:"deleted"                  ` //
	DeletedAt             *gtime.Time `json:"deletedAt"             orm:"deleted_at"               ` //
	UpdatedAt             *gtime.Time `json:"updatedAt"             orm:"updated_at"               ` //
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"               ` //
}
