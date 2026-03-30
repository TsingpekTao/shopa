// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewRecord is the golang structure for table review_record.
type ReviewRecord struct {
	Id               uint64      `json:"id"               orm:"id"                 description:""` //
	ReviewNo         string      `json:"reviewNo"         orm:"review_no"          description:""` //
	OrderNo          string      `json:"orderNo"          orm:"order_no"           description:""` //
	SubOrderNo       string      `json:"subOrderNo"       orm:"sub_order_no"       description:""` //
	ItemNo           string      `json:"itemNo"           orm:"item_no"            description:""` //
	UserId           uint64      `json:"userId"           orm:"user_id"            description:""` //
	ShopNo           string      `json:"shopNo"           orm:"shop_no"            description:""` //
	SpuNo            string      `json:"spuNo"            orm:"spu_no"             description:""` //
	SkuNo            string      `json:"skuNo"            orm:"sku_no"             description:""` //
	Score            uint        `json:"score"            orm:"score"              description:""` //
	Content          string      `json:"content"          orm:"content"            description:""` //
	MediasJson       string      `json:"mediasJson"       orm:"medias_json"        description:""` //
	Anonymous        int         `json:"anonymous"        orm:"anonymous"          description:""` //
	AppendContent    string      `json:"appendContent"    orm:"append_content"     description:""` //
	AppendMediasJson string      `json:"appendMediasJson" orm:"append_medias_json" description:""` //
	AppendAt         *gtime.Time `json:"appendAt"         orm:"append_at"          description:""` //
	SellerReply      string      `json:"sellerReply"      orm:"seller_reply"       description:""` //
	SellerReplyAt    *gtime.Time `json:"sellerReplyAt"    orm:"seller_reply_at"    description:""` //
	ReviewStatus     uint        `json:"reviewStatus"     orm:"review_status"      description:""` //
	LikeCount        uint64      `json:"likeCount"        orm:"like_count"         description:""` //
	Version          uint64      `json:"version"          orm:"version"            description:""` //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:""` //
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         description:""` //
	DeletedAt        *gtime.Time `json:"deletedAt"        orm:"deleted_at"         description:""` //
}
