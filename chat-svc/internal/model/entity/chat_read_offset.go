// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatReadOffset is the golang structure for table chat_read_offset.
type ChatReadOffset struct {
	Id              uint64      `json:"id"              orm:"id"                 description:""` //
	ConversationNo  string      `json:"conversationNo"  orm:"conversation_no"    description:""` //
	ReaderType      int         `json:"readerType"      orm:"reader_type"        description:""` //
	ReaderUserId    uint64      `json:"readerUserId"    orm:"reader_user_id"     description:""` //
	ReadToMessageNo string      `json:"readToMessageNo" orm:"read_to_message_no" description:""` //
	ReadToMessageId uint64      `json:"readToMessageId" orm:"read_to_message_id" description:""` //
	ReadAt          *gtime.Time `json:"readAt"          orm:"read_at"            description:""` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"         description:""` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"         description:""` //
}
