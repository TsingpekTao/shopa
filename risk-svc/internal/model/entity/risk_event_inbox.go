// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskEventInbox is the golang structure for table risk_event_inbox.
type RiskEventInbox struct {
	Id          uint64      `json:"id"          orm:"id"           ` //
	EventId     string      `json:"eventId"     orm:"event_id"     ` //
	EventType   string      `json:"eventType"   orm:"event_type"   ` //
	UserId      uint64      `json:"userId"      orm:"user_id"      ` //
	BizNo       string      `json:"bizNo"       orm:"biz_no"       ` //
	PayloadJson string      `json:"payloadJson" orm:"payload_json" ` //
	OccurredAt  *gtime.Time `json:"occurredAt"  orm:"occurred_at"  ` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` //
}
