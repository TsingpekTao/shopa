// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CartItemBackup is the golang structure of table cart_item_backup for DAO operations like Where/Data.
type CartItemBackup struct {
	g.Meta            `orm:"table:cart_item_backup, do:true"`
	Id                any         //
	UserId            any         //
	SkuNo             any         //
	SpuNo             any         //
	ShopNo            any         //
	Qty               any         //
	Checked           any         //
	Status            any         // 1=ACTIVE,2=INVALID
	InvalidReasonCode any         //
	SpuTitle          any         //
	SkuName           any         //
	SkuImageAssetId   any         //
	SalePrice         any         //
	MarketPrice       any         //
	SaleAttrsJson     any         //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}
