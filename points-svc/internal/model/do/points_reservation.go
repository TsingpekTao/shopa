// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsReservation is the golang structure of table points_reservation for DAO operations like Where/Data.
type PointsReservation struct {
	g.Meta                `orm:"table:points_reservation, do:true"`
	ReservationNo         any         //
	UserId                any         //
	OrderNo               any         //
	ShopNo                any         //
	ReservationStatusCode any         // LOCKED/CONFIRMED/CANCELED/EXPIRED
	RequestedPoints       any         //
	LockedPoints          any         //
	LockedCashAmountCent  any         //
	DeductionDigest       any         //
	RuleSnapshotJson      any         //
	IdempotencyKey        any         //
	ExpireAt              *gtime.Time //
	ConfirmedAt           *gtime.Time //
	CanceledAt            *gtime.Time //
	CancelReasonCode      any         //
	CreatedAt             *gtime.Time //
	UpdatedAt             *gtime.Time //
}
