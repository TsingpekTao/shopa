// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryReservation is the golang structure for table inventory_reservation.
type InventoryReservation struct {
	Id                uint64      `json:"id"                orm:"id"                 ` //
	ReservationNo     string      `json:"reservationNo"     orm:"reservation_no"     ` //
	OrderNo           string      `json:"orderNo"           orm:"order_no"           ` //
	UserId            uint64      `json:"userId"            orm:"user_id"            ` //
	ReservationStatus uint        `json:"reservationStatus" orm:"reservation_status" ` //
	ReserveMode       uint        `json:"reserveMode"       orm:"reserve_mode"       ` //
	ExpiredAt         *gtime.Time `json:"expiredAt"         orm:"expired_at"         ` //
	ConfirmedAt       *gtime.Time `json:"confirmedAt"       orm:"confirmed_at"       ` //
	CanceledAt        *gtime.Time `json:"canceledAt"        orm:"canceled_at"        ` //
	CancelReasonCode  string      `json:"cancelReasonCode"  orm:"cancel_reason_code" ` //
	IdempotencyKey    string      `json:"idempotencyKey"    orm:"idempotency_key"    ` //
	RequestId         string      `json:"requestId"         orm:"request_id"         ` //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"         ` //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"         ` //
}
