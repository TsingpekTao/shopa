// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewIdempotency is the golang structure of table review_idempotency for DAO operations like Where/Data.
type ReviewIdempotency struct {
	g.Meta         `orm:"table:review_idempotency, do:true"`
	Id             any         //
	UserId         any         //
	Action         any         //
	IdempotencyKey any         //
	ResourceNo     any         //
	Status         any         //
	ResponseJson   any         //
	ErrorCode      any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
