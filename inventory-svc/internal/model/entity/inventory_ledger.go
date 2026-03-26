// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryLedger 是 inventory_ledger 表的 Go 结构体。
type InventoryLedger struct {
	Id                uint64      `json:"id"                orm:"id"                  ` //
	TxnNo             string      `json:"txnNo"             orm:"txn_no"              ` //
	BizType           string      `json:"bizType"           orm:"biz_type"            ` //
	BizNo             string      `json:"bizNo"             orm:"biz_no"              ` //
	ActionCode        string      `json:"actionCode"        orm:"action_code"         ` //
	ReservationNo     string      `json:"reservationNo"     orm:"reservation_no"      ` //
	OrderNo           string      `json:"orderNo"           orm:"order_no"            ` //
	SkuNo             string      `json:"skuNo"             orm:"sku_no"              ` //
	SpuNo             string      `json:"spuNo"             orm:"spu_no"              ` //
	ShopNo            string      `json:"shopNo"            orm:"shop_no"             ` //
	DeltaTotalQty     int64       `json:"deltaTotalQty"     orm:"delta_total_qty"     ` //
	DeltaLockedQty    int64       `json:"deltaLockedQty"    orm:"delta_locked_qty"    ` //
	DeltaAvailableQty int64       `json:"deltaAvailableQty" orm:"delta_available_qty" ` //
	AfterTotalQty     uint64      `json:"afterTotalQty"     orm:"after_total_qty"     ` //
	AfterLockedQty    uint64      `json:"afterLockedQty"    orm:"after_locked_qty"    ` //
	AfterAvailableQty uint64      `json:"afterAvailableQty" orm:"after_available_qty" ` //
	StockVersion      uint64      `json:"stockVersion"      orm:"stock_version"       ` //
	OperatorType      string      `json:"operatorType"      orm:"operator_type"       ` // SELLER/ORDER/ADMIN/SYSTEM
	OperatorUserId    uint64      `json:"operatorUserId"    orm:"operator_user_id"    ` //
	RequestId         string      `json:"requestId"         orm:"request_id"          ` //
	Remark            string      `json:"remark"            orm:"remark"              ` //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          ` //
}
