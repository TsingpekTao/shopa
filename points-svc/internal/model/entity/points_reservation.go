// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsReservation is the golang structure for table points_reservation.
type PointsReservation struct {
	ReservationNo         string      `json:"reservationNo"         orm:"reservation_no"          ` //
	UserId                uint64      `json:"userId"                orm:"user_id"                 ` //
	OrderNo               string      `json:"orderNo"               orm:"order_no"                ` //
	ShopNo                string      `json:"shopNo"                orm:"shop_no"                 ` //
	ReservationStatusCode string      `json:"reservationStatusCode" orm:"reservation_status_code" ` // LOCKED/CONFIRMED/CANCELED/EXPIRED
	RequestedPoints       uint64      `json:"requestedPoints"       orm:"requested_points"        ` //
	LockedPoints          uint64      `json:"lockedPoints"          orm:"locked_points"           ` //
	LockedCashAmountCent  int64       `json:"lockedCashAmountCent"  orm:"locked_cash_amount_cent" ` //
	DeductionDigest       string      `json:"deductionDigest"       orm:"deduction_digest"        ` //
	RuleSnapshotJson      string      `json:"ruleSnapshotJson"      orm:"rule_snapshot_json"      ` //
	IdempotencyKey        string      `json:"idempotencyKey"        orm:"idempotency_key"         ` //
	ExpireAt              *gtime.Time `json:"expireAt"              orm:"expire_at"               ` //
	ConfirmedAt           *gtime.Time `json:"confirmedAt"           orm:"confirmed_at"            ` //
	CanceledAt            *gtime.Time `json:"canceledAt"            orm:"canceled_at"             ` //
	CancelReasonCode      string      `json:"cancelReasonCode"      orm:"cancel_reason_code"      ` //
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"              ` //
	UpdatedAt             *gtime.Time `json:"updatedAt"             orm:"updated_at"              ` //
}
