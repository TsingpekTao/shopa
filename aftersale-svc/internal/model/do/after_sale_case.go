// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSaleCase is the golang structure of table after_sale_case for DAO operations like Where/Data.
type AfterSaleCase struct {
	g.Meta               `orm:"table:after_sale_case, do:true"`
	Id                   any         //
	AfterSaleNo          any         //
	OrderNo              any         //
	SubOrderNo           any         //
	ItemNo               any         //
	UserId               any         //
	ShopNo               any         //
	SpuNo                any         //
	SkuNo                any         //
	Qty                  any         //
	AfterSaleType        any         //
	AfterSaleStatus      any         //
	ApplyRefundAmount    any         //
	ApprovedRefundAmount any         //
	ReasonCode           any         //
	ReasonDesc           any         //
	EvidenceAssetIdsJson any         //
	BuyerRemark          any         //
	SellerReply          any         //
	RejectReasonCode     any         //
	CancelReasonCode     any         //
	Version              any         //
	ClosedAt             *gtime.Time //
	CreatedAt            *gtime.Time //
	UpdatedAt            *gtime.Time //
	DeletedAt            *gtime.Time //
}
