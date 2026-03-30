// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderSub is the golang structure for table order_sub.
type OrderSub struct {
	Id             uint64      `json:"id"             orm:"id"              description:""` //
	SubOrderNo     string      `json:"subOrderNo"     orm:"sub_order_no"    description:""` //
	OrderNo        string      `json:"orderNo"        orm:"order_no"        description:""` //
	ShopNo         string      `json:"shopNo"         orm:"shop_no"         description:""` //
	SubStatus      uint        `json:"subStatus"      orm:"sub_status"      description:""` //
	GoodsAmount    uint64      `json:"goodsAmount"    orm:"goods_amount"    description:""` //
	FreightAmount  uint64      `json:"freightAmount"  orm:"freight_amount"  description:""` //
	DiscountAmount uint64      `json:"discountAmount" orm:"discount_amount" description:""` //
	PayableAmount  uint64      `json:"payableAmount"  orm:"payable_amount"  description:""` //
	PaidAmount     uint64      `json:"paidAmount"     orm:"paid_amount"     description:""` //
	SellerRemark   string      `json:"sellerRemark"   orm:"seller_remark"   description:""` //
	BuyerRemark    string      `json:"buyerRemark"    orm:"buyer_remark"    description:""` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""` //
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"      description:""` //
}
