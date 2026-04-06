// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentFeedback is the golang structure of table agent_feedback for DAO operations like Where/Data.
type AgentFeedback struct {
	g.Meta         `orm:"table:agent_feedback, do:true"`
	Id             any         //
	ConversationNo any         //
	RunNo          any         //
	UserId         any         //
	ShopNo         any         //
	FeedbackCode   any         //
	Resolved       any         //
	Comment        any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
