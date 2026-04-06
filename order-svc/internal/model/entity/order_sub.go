// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderSub is the golang structure for table order_sub.
type OrderSub struct {
	Id                   uint64      `json:"id"                   orm:"id"                     ` //
	SubOrderNo           string      `json:"subOrderNo"           orm:"sub_order_no"           ` //
	OrderNo              string      `json:"orderNo"              orm:"order_no"               ` //
	ShopNo               string      `json:"shopNo"               orm:"shop_no"                ` //
	SubStatus            uint        `json:"subStatus"            orm:"sub_status"             ` //
	GoodsAmount          uint64      `json:"goodsAmount"          orm:"goods_amount"           ` //
	FreightAmount        uint64      `json:"freightAmount"        orm:"freight_amount"         ` //
	DiscountAmount       uint64      `json:"discountAmount"       orm:"discount_amount"        ` //
	PayableAmount        uint64      `json:"payableAmount"        orm:"payable_amount"         ` //
	PaidAmount           uint64      `json:"paidAmount"           orm:"paid_amount"            ` //
	PointsUsed           uint64      `json:"pointsUsed"           orm:"points_used"            ` //
	PointsDiscountAmount uint64      `json:"pointsDiscountAmount" orm:"points_discount_amount" ` //
	SellerRemark         string      `json:"sellerRemark"         orm:"seller_remark"          ` //
	BuyerRemark          string      `json:"buyerRemark"          orm:"buyer_remark"           ` //
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             ` //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             ` //
	DeletedAt            *gtime.Time `json:"deletedAt"            orm:"deleted_at"             ` //
}
