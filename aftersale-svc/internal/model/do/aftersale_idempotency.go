// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AftersaleIdempotency is the golang structure of table aftersale_idempotency for DAO operations like Where/Data.
type AftersaleIdempotency struct {
	g.Meta         `orm:"table:aftersale_idempotency, do:true"`
	Id             any         //
	UserId         any         //
	IdempotencyKey any         //
	ActionCode     any         //
	ResourceNo     any         //
	Status         any         //
	ResponseJson   any         //
	ErrorCode      any         //
	ExpireAt       *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
