// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentMessage is the golang structure of table agent_message for DAO operations like Where/Data.
type AgentMessage struct {
	g.Meta          `orm:"table:agent_message, do:true"`
	Id              any         //
	MessageNo       any         //
	ConversationNo  any         //
	RunNo           any         //
	ReplyToTurnNo   any         //
	SenderTypeCode  any         //
	MessageTypeCode any         //
	ClientMessageNo any         //
	ContentText     any         //
	AssetIdsJson    any         //
	ExtJson         any         //
	Interrupted     any         //
	SentAt          *gtime.Time //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       *gtime.Time //
}
