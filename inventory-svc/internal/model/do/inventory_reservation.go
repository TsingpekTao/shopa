// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryReservation 是用于 DAO 操作（如 Where/Data）的 inventory_reservation 表 Go 结构体。
type InventoryReservation struct {
	g.Meta            `orm:"table:inventory_reservation, do:true"`
	Id                any         //
	ReservationNo     any         //
	OrderNo           any         //
	UserId            any         //
	ReservationStatus any         //
	ReserveMode       any         //
	ExpiredAt         *gtime.Time //
	ConfirmedAt       *gtime.Time //
	CanceledAt        *gtime.Time //
	CancelReasonCode  any         //
	IdempotencyKey    any         //
	RequestId         any         //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}
