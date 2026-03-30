// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatMessage is the golang structure for table chat_message.
type ChatMessage struct {
	Id              uint64      `json:"id"              orm:"id"                description:""` //
	MessageNo       string      `json:"messageNo"       orm:"message_no"        description:""` //
	ConversationNo  string      `json:"conversationNo"  orm:"conversation_no"   description:""` //
	SenderType      int         `json:"senderType"      orm:"sender_type"       description:""` //
	SenderUserId    uint64      `json:"senderUserId"    orm:"sender_user_id"    description:""` //
	ClientMessageNo string      `json:"clientMessageNo" orm:"client_message_no" description:""` //
	MessageType     int         `json:"messageType"     orm:"message_type"      description:""` //
	ContentText     string      `json:"contentText"     orm:"content_text"      description:""` //
	MediaAssetId    uint64      `json:"mediaAssetId"    orm:"media_asset_id"    description:""` //
	ExtJson         string      `json:"extJson"         orm:"ext_json"          description:""` //
	SentAt          *gtime.Time `json:"sentAt"          orm:"sent_at"           description:""` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        description:""` //
}
