// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryIdempotency is the golang structure of table inventory_idempotency for DAO operations like Where/Data.
type InventoryIdempotency struct {
	g.Meta         `orm:"table:inventory_idempotency, do:true"`
	Id             any         //
	ActorScope     any         // SELLER/ORDER/ADMIN/SYSTEM
	ActorId        any         //
	ActionCode     any         //
	IdempotencyKey any         //
	RequestHash    any         //
	Status         any         //
	ResponseCode   any         //
	ResponseJson   any         //
	ExpiredAt      *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
