// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentIdempotency is the golang structure of table fulfillment_idempotency for DAO operations like Where/Data.
type FulfillmentIdempotency struct {
	g.Meta         `orm:"table:fulfillment_idempotency, do:true"`
	Id             any         //
	IdempotencyKey any         //
	Scope          any         //
	RequestHash    any         //
	ResourceId     any         //
	ResponseJson   any         //
	Status         any         // 1 SUCCEEDED 2 PROCESSING 3 FAILED
	ExpireAt       *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
