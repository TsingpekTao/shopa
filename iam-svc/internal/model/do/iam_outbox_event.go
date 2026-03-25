// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamOutboxEvent is the golang structure of table iam_outbox_event for DAO operations like Where/Data.
type IamOutboxEvent struct {
	g.Meta      `orm:"table:iam_outbox_event, do:true"`
	Id          any         //
	EventId     any         // Global unique event id
	EventType   any         //
	PayloadJson any         //
	Status      any         // 1 NEW,2 PROCESSING,3 SENT,4 FAILED,5 DLQ
	AvailableAt *gtime.Time // Next dispatch schedule time
	SentAt      *gtime.Time // Successful dispatch time
	FailCount   any         // Dispatch failure count
	LastError   any         // Last dispatch error text
	RetryCount  any         //
	NextRetryAt *gtime.Time //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
