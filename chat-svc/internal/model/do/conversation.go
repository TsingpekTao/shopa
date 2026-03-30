// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Conversation is the golang structure of table conversation for DAO operations like Where/Data.
type Conversation struct {
	g.Meta             `orm:"table:conversation, do:true"`
	Id                 any         //
	ConversationNo     any         //
	BuyerId            any         //
	ShopNo             any         //
	SceneCode          any         //
	OrderNo            any         //
	SubOrderNo         any         //
	AnchorSpuNo        any         //
	AnchorSkuNo        any         //
	LastMessageNo      any         //
	LastMessagePreview any         //
	LastMessageAt      *gtime.Time //
	ConversationStatus any         //
	Version            any         //
	CreatedAt          *gtime.Time //
	UpdatedAt          *gtime.Time //
	DeletedAt          *gtime.Time //
}
