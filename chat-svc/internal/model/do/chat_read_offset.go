// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatReadOffset is the golang structure of table chat_read_offset for DAO operations like Where/Data.
type ChatReadOffset struct {
	g.Meta          `orm:"table:chat_read_offset, do:true"`
	Id              any         //
	ConversationNo  any         //
	ReaderType      any         //
	ReaderUserId    any         //
	ReadToMessageNo any         //
	ReadToMessageId any         //
	ReadAt          *gtime.Time //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
