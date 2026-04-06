// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentFeedback is the golang structure for table agent_feedback.
type AgentFeedback struct {
	Id             uint64      `json:"id"             orm:"id"              ` //
	ConversationNo string      `json:"conversationNo" orm:"conversation_no" ` //
	RunNo          string      `json:"runNo"          orm:"run_no"          ` //
	UserId         uint64      `json:"userId"         orm:"user_id"         ` //
	ShopNo         string      `json:"shopNo"         orm:"shop_no"         ` //
	FeedbackCode   string      `json:"feedbackCode"   orm:"feedback_code"   ` //
	Resolved       int         `json:"resolved"       orm:"resolved"        ` //
	Comment        string      `json:"comment"        orm:"comment"         ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"      ` //
}
