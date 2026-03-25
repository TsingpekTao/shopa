// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogConsumerEventDedup is the golang structure for table catalog_consumer_event_dedup.
type CatalogConsumerEventDedup struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	ConsumerName string      `json:"consumerName" orm:"consumer_name" ` //
	EventId      string      `json:"eventId"      orm:"event_id"      ` //
	PayloadHash  string      `json:"payloadHash"  orm:"payload_hash"  ` //
	FirstSeenAt  *gtime.Time `json:"firstSeenAt"  orm:"first_seen_at" ` //
}
