// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentMessage is the golang structure for table agent_message.
type AgentMessage struct {
	Id              uint64      `json:"id"              orm:"id"                ` //
	MessageNo       string      `json:"messageNo"       orm:"message_no"        ` //
	ConversationNo  string      `json:"conversationNo"  orm:"conversation_no"   ` //
	RunNo           string      `json:"runNo"           orm:"run_no"            ` //
	ReplyToTurnNo   uint64      `json:"replyToTurnNo"   orm:"reply_to_turn_no"  ` //
	SenderTypeCode  string      `json:"senderTypeCode"  orm:"sender_type_code"  ` //
	MessageTypeCode string      `json:"messageTypeCode" orm:"message_type_code" ` //
	ClientMessageNo string      `json:"clientMessageNo" orm:"client_message_no" ` //
	ContentText     string      `json:"contentText"     orm:"content_text"      ` //
	AssetIdsJson    string      `json:"assetIdsJson"    orm:"asset_ids_json"    ` //
	ExtJson         string      `json:"extJson"         orm:"ext_json"          ` //
	Interrupted     int         `json:"interrupted"     orm:"interrupted"       ` //
	SentAt          *gtime.Time `json:"sentAt"          orm:"sent_at"           ` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"        ` //
	DeletedAt       *gtime.Time `json:"deletedAt"       orm:"deleted_at"        ` //
}
