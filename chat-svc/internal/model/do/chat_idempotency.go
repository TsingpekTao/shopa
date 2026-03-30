// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatIdempotency is the golang structure of table chat_idempotency for DAO operations like Where/Data.
type ChatIdempotency struct {
	g.Meta         `orm:"table:chat_idempotency, do:true"`
	Id             any         //
	UserId         any         //
	BizCode        any         //
	IdempotencyKey any         //
	TargetNo       any         //
	ResponseJson   any         //
	Status         any         //
	ExpiredAt      *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
