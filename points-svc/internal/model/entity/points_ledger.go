// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsLedger is the golang structure for table points_ledger.
type PointsLedger struct {
	Id              uint64      `json:"id"              orm:"id"                ` //
	LedgerNo        string      `json:"ledgerNo"        orm:"ledger_no"         ` // Business ledger identifier
	UserId          uint64      `json:"userId"          orm:"user_id"           ` //
	EntryTypeCode   string      `json:"entryTypeCode"   orm:"entry_type_code"   ` // INIT/LOCK/CONFIRM/CANCEL/GRANT/RETURN/REVERSE/EXPIRE/ADJUST/FREEZE/UNFREEZE
	BizType         string      `json:"bizType"         orm:"biz_type"          ` // REGISTER_INIT/ORDER_PAY/REFUND/etc
	BizNo           string      `json:"bizNo"           orm:"biz_no"            ` // Business identifier for idempotency
	ReservationNo   string      `json:"reservationNo"   orm:"reservation_no"    ` // Reservation identifier if relevant
	RelatedBucketNo string      `json:"relatedBucketNo" orm:"related_bucket_no" ` // Bucket identifier if relevant
	PointsDelta     int64       `json:"pointsDelta"     orm:"points_delta"      ` // Points delta for this ledger entry
	AvailableAfter  int64       `json:"availableAfter"  orm:"available_after"   ` // Available balance snapshot after apply
	FrozenAfter     int64       `json:"frozenAfter"     orm:"frozen_after"      ` // Frozen balance snapshot after apply
	DebtAfter       uint64      `json:"debtAfter"       orm:"debt_after"        ` // Debt snapshot after apply
	CashAmountCent  int64       `json:"cashAmountCent"  orm:"cash_amount_cent"  ` // Related cash amount in cents
	Remark          string      `json:"remark"          orm:"remark"            ` // Operator remark or domain explanation
	ExtraJson       string      `json:"extraJson"       orm:"extra_json"        ` // Extended metadata snapshot
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        ` //
}
