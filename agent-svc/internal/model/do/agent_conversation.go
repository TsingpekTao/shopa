// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentConversation is the golang structure of table agent_conversation for DAO operations like Where/Data.
type AgentConversation struct {
	g.Meta                  `orm:"table:agent_conversation, do:true"`
	Id                      any         //
	ConversationNo          any         //
	UserId                  any         //
	BotCode                 any         //
	SceneCode               any         //
	ShopNo                  any         //
	OrderNo                 any         //
	SubOrderNo              any         //
	AnchorSpuNo             any         //
	AnchorSkuNo             any         //
	ConversationStatusCode  any         //
	LastRunStatusCode       any         //
	IsHumanHandover         any         //
	PendingMessageCount     any         //
	SessionSummary          any         //
	SessionSummaryVersion   any         //
	LastSummarizedMessageNo any         //
	LastMessageAt           *gtime.Time //
	LastQueueNoticeAt       *gtime.Time //
	CreatedAt               *gtime.Time //
	UpdatedAt               *gtime.Time //
	DeletedAt               *gtime.Time //
}
