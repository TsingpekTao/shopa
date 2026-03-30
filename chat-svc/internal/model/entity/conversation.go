// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Conversation is the golang structure for table conversation.
type Conversation struct {
	Id                 uint64      `json:"id"                 orm:"id"                   description:""` //
	ConversationNo     string      `json:"conversationNo"     orm:"conversation_no"      description:""` //
	BuyerId            uint64      `json:"buyerId"            orm:"buyer_id"             description:""` //
	ShopNo             string      `json:"shopNo"             orm:"shop_no"              description:""` //
	SceneCode          string      `json:"sceneCode"          orm:"scene_code"           description:""` //
	OrderNo            string      `json:"orderNo"            orm:"order_no"             description:""` //
	SubOrderNo         string      `json:"subOrderNo"         orm:"sub_order_no"         description:""` //
	AnchorSpuNo        string      `json:"anchorSpuNo"        orm:"anchor_spu_no"        description:""` //
	AnchorSkuNo        string      `json:"anchorSkuNo"        orm:"anchor_sku_no"        description:""` //
	LastMessageNo      string      `json:"lastMessageNo"      orm:"last_message_no"      description:""` //
	LastMessagePreview string      `json:"lastMessagePreview" orm:"last_message_preview" description:""` //
	LastMessageAt      *gtime.Time `json:"lastMessageAt"      orm:"last_message_at"      description:""` //
	ConversationStatus int         `json:"conversationStatus" orm:"conversation_status"  description:""` //
	Version            uint64      `json:"version"            orm:"version"              description:""` //
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:""` //
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:""` //
	DeletedAt          *gtime.Time `json:"deletedAt"          orm:"deleted_at"           description:""` //
}
