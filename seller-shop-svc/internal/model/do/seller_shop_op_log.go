// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerShopOpLog is the golang structure of table seller_shop_op_log for DAO operations like Where/Data.
type SellerShopOpLog struct {
	g.Meta         `orm:"table:seller_shop_op_log, do:true"`
	Id             any         //
	ShopNo         any         //
	OwnerUserId    any         //
	OperatorUserId any         //
	ActionCode     any         // PROVISION_START/ACTIVATE/FREEZE/CLOSE/REOPEN
	FromStatus     any         //
	ToStatus       any         //
	ReasonCode     any         //
	Reason         any         //
	RequestId      any         //
	EventId        any         //
	CreatedAt      *gtime.Time //
}
