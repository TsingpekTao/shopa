// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderIdempotency is the golang structure of table order_idempotency for DAO operations like Where/Data.
type OrderIdempotency struct {
	g.Meta         `orm:"table:order_idempotency, do:true"`
	Id             any         //
	UserId         any         //
	IdempotencyKey any         //
	ActionCode     any         //
	OrderNo        any         //
	Status         any         //
	ResponseJson   any         //
	ErrorCode      any         //
	ExpireAt       *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
