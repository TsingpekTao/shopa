// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSpu is the golang structure for table catalog_spu.
type CatalogSpu struct {
	Id                      uint64      `json:"id"                      orm:"id"                          ` //
	SpuNo                   string      `json:"spuNo"                   orm:"spu_no"                      ` //
	ShopNo                  string      `json:"shopNo"                  orm:"shop_no"                     ` //
	Title                   string      `json:"title"                   orm:"title"                       ` //
	SubTitle                string      `json:"subTitle"                orm:"sub_title"                   ` //
	CategoryId              uint64      `json:"categoryId"              orm:"category_id"                 ` //
	BrandNo                 string      `json:"brandNo"                 orm:"brand_no"                    ` //
	MainImageAssetIdsJson   string      `json:"mainImageAssetIdsJson"   orm:"main_image_asset_ids_json"   ` //
	DetailImageAssetIdsJson string      `json:"detailImageAssetIdsJson" orm:"detail_image_asset_ids_json" ` //
	SpuStatus               uint        `json:"spuStatus"               orm:"spu_status"                  ` //
	SpuStockStatus          uint        `json:"spuStockStatus"          orm:"spu_stock_status"            ` //
	MinSalePrice            uint64      `json:"minSalePrice"            orm:"min_sale_price"              ` //
	MaxSalePrice            uint64      `json:"maxSalePrice"            orm:"max_sale_price"              ` //
	MinMarketPrice          uint64      `json:"minMarketPrice"          orm:"min_market_price"            ` //
	MaxMarketPrice          uint64      `json:"maxMarketPrice"          orm:"max_market_price"            ` //
	PublishTime             *gtime.Time `json:"publishTime"             orm:"publish_time"                ` //
	Version                 uint        `json:"version"                 orm:"version"                     ` //
	ReviewStatus            uint        `json:"reviewStatus"            orm:"review_status"               ` // 1 PENDING, 2 APPROVED, 3 REJECTED
	RejectReasonCode        string      `json:"rejectReasonCode"        orm:"reject_reason_code"          ` //
	RejectComment           string      `json:"rejectComment"           orm:"reject_comment"              ` //
	ReviewComment           string      `json:"reviewComment"           orm:"review_comment"              ` //
	ReviewerId              uint64      `json:"reviewerId"              orm:"reviewer_id"                 ` //
	SubmittedAt             *gtime.Time `json:"submittedAt"             orm:"submitted_at"                ` //
	ReviewedAt              *gtime.Time `json:"reviewedAt"              orm:"reviewed_at"                 ` //
	SoldCount               uint64      `json:"soldCount"               orm:"sold_count"                  ` //
	DeletedAt               *gtime.Time `json:"deletedAt"               orm:"deleted_at"                  ` //
	CreatedAt               *gtime.Time `json:"createdAt"               orm:"created_at"                  ` //
	UpdatedAt               *gtime.Time `json:"updatedAt"               orm:"updated_at"                  ` //
}
