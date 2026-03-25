// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaConsumerEventDedup is the golang structure for table media_consumer_event_dedup.
type MediaConsumerEventDedup struct {
	Id           uint64      `json:"id"           orm:"id"            description:""` //
	ConsumerName string      `json:"consumerName" orm:"consumer_name" description:""` //
	EventId      string      `json:"eventId"      orm:"event_id"      description:""` //
	PayloadHash  string      `json:"payloadHash"  orm:"payload_hash"  description:""` //
	FirstSeenAt  *gtime.Time `json:"firstSeenAt"  orm:"first_seen_at" description:""` //
}
