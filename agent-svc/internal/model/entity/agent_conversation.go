// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentConversation is the golang structure for table agent_conversation.
type AgentConversation struct {
	Id                      uint64      `json:"id"                      orm:"id"                         ` //
	ConversationNo          string      `json:"conversationNo"          orm:"conversation_no"            ` //
	UserId                  uint64      `json:"userId"                  orm:"user_id"                    ` //
	BotCode                 string      `json:"botCode"                 orm:"bot_code"                   ` //
	SceneCode               string      `json:"sceneCode"               orm:"scene_code"                 ` //
	ShopNo                  string      `json:"shopNo"                  orm:"shop_no"                    ` //
	OrderNo                 string      `json:"orderNo"                 orm:"order_no"                   ` //
	SubOrderNo              string      `json:"subOrderNo"              orm:"sub_order_no"               ` //
	AnchorSpuNo             string      `json:"anchorSpuNo"             orm:"anchor_spu_no"              ` //
	AnchorSkuNo             string      `json:"anchorSkuNo"             orm:"anchor_sku_no"              ` //
	ConversationStatusCode  string      `json:"conversationStatusCode"  orm:"conversation_status_code"   ` //
	LastRunStatusCode       string      `json:"lastRunStatusCode"       orm:"last_run_status_code"       ` //
	IsHumanHandover         int         `json:"isHumanHandover"         orm:"is_human_handover"          ` //
	PendingMessageCount     uint64      `json:"pendingMessageCount"     orm:"pending_message_count"      ` //
	SessionSummary          string      `json:"sessionSummary"          orm:"session_summary"            ` //
	SessionSummaryVersion   uint64      `json:"sessionSummaryVersion"   orm:"session_summary_version"    ` //
	LastSummarizedMessageNo string      `json:"lastSummarizedMessageNo" orm:"last_summarized_message_no" ` //
	LastMessageAt           *gtime.Time `json:"lastMessageAt"           orm:"last_message_at"            ` //
	LastQueueNoticeAt       *gtime.Time `json:"lastQueueNoticeAt"       orm:"last_queue_notice_at"       ` //
	CreatedAt               *gtime.Time `json:"createdAt"               orm:"created_at"                 ` //
	UpdatedAt               *gtime.Time `json:"updatedAt"               orm:"updated_at"                 ` //
	DeletedAt               *gtime.Time `json:"deletedAt"               orm:"deleted_at"                 ` //
}
