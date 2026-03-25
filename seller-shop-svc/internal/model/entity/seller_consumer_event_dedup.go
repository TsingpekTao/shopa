// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerConsumerEventDedup is the golang structure for table seller_consumer_event_dedup.
type SellerConsumerEventDedup struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	ConsumerName string      `json:"consumerName" orm:"consumer_name" ` // seller-role-assigned-consumer/...
	EventId      string      `json:"eventId"      orm:"event_id"      ` // Incoming event id
	PayloadHash  string      `json:"payloadHash"  orm:"payload_hash"  ` //
	FirstSeenAt  *gtime.Time `json:"firstSeenAt"  orm:"first_seen_at" ` //
}
