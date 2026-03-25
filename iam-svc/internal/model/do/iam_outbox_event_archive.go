// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamOutboxEventArchive is the golang structure of table iam_outbox_event_archive for DAO operations like Where/Data.
type IamOutboxEventArchive struct {
	g.Meta      `orm:"table:iam_outbox_event_archive, do:true"`
	Id          any         //
	EventId     any         //
	EventType   any         //
	PayloadJson any         //
	Status      any         // Archive rows are usually SENT/FAILED/DLQ snapshots
	AvailableAt *gtime.Time //
	SentAt      *gtime.Time //
	FailCount   any         //
	LastError   any         //
	RetryCount  any         //
	NextRetryAt *gtime.Time //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	ArchivedAt  *gtime.Time //
}
