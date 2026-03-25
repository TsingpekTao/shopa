// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSpu is the golang structure of table catalog_spu for DAO operations like Where/Data.
type CatalogSpu struct {
	g.Meta                  `orm:"table:catalog_spu, do:true"`
	Id                      any         //
	SpuNo                   any         //
	ShopNo                  any         //
	Title                   any         //
	SubTitle                any         //
	CategoryId              any         //
	BrandNo                 any         //
	MainImageAssetIdsJson   any         //
	DetailImageAssetIdsJson any         //
	SpuStatus               any         //
	SpuStockStatus          any         //
	MinSalePrice            any         //
	MaxSalePrice            any         //
	MinMarketPrice          any         //
	MaxMarketPrice          any         //
	PublishTime             *gtime.Time //
	Version                 any         //
	ReviewStatus            any         // 1 PENDING, 2 APPROVED, 3 REJECTED
	RejectReasonCode        any         //
	RejectComment           any         //
	ReviewComment           any         //
	ReviewerId              any         //
	SubmittedAt             *gtime.Time //
	ReviewedAt              *gtime.Time //
	SoldCount               any         //
	DeletedAt               *gtime.Time //
	CreatedAt               *gtime.Time //
	UpdatedAt               *gtime.Time //
}
