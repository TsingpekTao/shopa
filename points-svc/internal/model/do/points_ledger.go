// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsLedger is the golang structure of table points_ledger for DAO operations like Where/Data.
type PointsLedger struct {
	g.Meta          `orm:"table:points_ledger, do:true"`
	Id              any         //
	LedgerNo        any         // Business ledger identifier
	UserId          any         //
	EntryTypeCode   any         // INIT/LOCK/CONFIRM/CANCEL/GRANT/RETURN/REVERSE/EXPIRE/ADJUST/FREEZE/UNFREEZE
	BizType         any         // REGISTER_INIT/ORDER_PAY/REFUND/etc
	BizNo           any         // Business identifier for idempotency
	ReservationNo   any         // Reservation identifier if relevant
	RelatedBucketNo any         // Bucket identifier if relevant
	PointsDelta     any         // Points delta for this ledger entry
	AvailableAfter  any         // Available balance snapshot after apply
	FrozenAfter     any         // Frozen balance snapshot after apply
	DebtAfter       any         // Debt snapshot after apply
	CashAmountCent  any         // Related cash amount in cents
	Remark          any         // Operator remark or domain explanation
	ExtraJson       any         // Extended metadata snapshot
	CreatedAt       *gtime.Time //
}
