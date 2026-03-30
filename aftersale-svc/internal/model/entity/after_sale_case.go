// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSaleCase is the golang structure for table after_sale_case.
type AfterSaleCase struct {
	Id                   uint64      `json:"id"                   orm:"id"                      ` //
	AfterSaleNo          string      `json:"afterSaleNo"          orm:"after_sale_no"           ` //
	OrderNo              string      `json:"orderNo"              orm:"order_no"                ` //
	SubOrderNo           string      `json:"subOrderNo"           orm:"sub_order_no"            ` //
	ItemNo               string      `json:"itemNo"               orm:"item_no"                 ` //
	UserId               uint64      `json:"userId"               orm:"user_id"                 ` //
	ShopNo               string      `json:"shopNo"               orm:"shop_no"                 ` //
	SpuNo                string      `json:"spuNo"                orm:"spu_no"                  ` //
	SkuNo                string      `json:"skuNo"                orm:"sku_no"                  ` //
	Qty                  uint        `json:"qty"                  orm:"qty"                     ` //
	AfterSaleType        uint        `json:"afterSaleType"        orm:"after_sale_type"         ` //
	AfterSaleStatus      uint        `json:"afterSaleStatus"      orm:"after_sale_status"       ` //
	ApplyRefundAmount    uint64      `json:"applyRefundAmount"    orm:"apply_refund_amount"     ` //
	ApprovedRefundAmount uint64      `json:"approvedRefundAmount" orm:"approved_refund_amount"  ` //
	ReasonCode           string      `json:"reasonCode"           orm:"reason_code"             ` //
	ReasonDesc           string      `json:"reasonDesc"           orm:"reason_desc"             ` //
	EvidenceAssetIdsJson string      `json:"evidenceAssetIdsJson" orm:"evidence_asset_ids_json" ` //
	BuyerRemark          string      `json:"buyerRemark"          orm:"buyer_remark"            ` //
	SellerReply          string      `json:"sellerReply"          orm:"seller_reply"            ` //
	RejectReasonCode     uint        `json:"rejectReasonCode"     orm:"reject_reason_code"      ` //
	CancelReasonCode     string      `json:"cancelReasonCode"     orm:"cancel_reason_code"      ` //
	Version              uint64      `json:"version"              orm:"version"                 ` //
	ClosedAt             *gtime.Time `json:"closedAt"             orm:"closed_at"               ` //
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"              ` //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"              ` //
	DeletedAt            *gtime.Time `json:"deletedAt"            orm:"deleted_at"              ` //
}
