// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryLedger 是用于 DAO 操作（如 Where/Data）的 inventory_ledger 表 Go 结构体。
type InventoryLedger struct {
	g.Meta            `orm:"table:inventory_ledger, do:true"`
	Id                any         //
	TxnNo             any         //
	BizType           any         //
	BizNo             any         //
	ActionCode        any         //
	ReservationNo     any         //
	OrderNo           any         //
	SkuNo             any         //
	SpuNo             any         //
	ShopNo            any         //
	DeltaTotalQty     any         //
	DeltaLockedQty    any         //
	DeltaAvailableQty any         //
	AfterTotalQty     any         //
	AfterLockedQty    any         //
	AfterAvailableQty any         //
	StockVersion      any         //
	OperatorType      any         // SELLER/ORDER/ADMIN/SYSTEM
	OperatorUserId    any         //
	RequestId         any         //
	Remark            any         //
	CreatedAt         *gtime.Time //
}
