// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventorySkuContext 是用于 DAO 操作（如 Where/Data）的 inventory_sku_context 表 Go 结构体。
type InventorySkuContext struct {
	g.Meta    `orm:"table:inventory_sku_context, do:true"`
	Id        any         //
	SkuNo     any         //
	SpuNo     any         //
	ShopNo    any         //
	Enabled   any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
