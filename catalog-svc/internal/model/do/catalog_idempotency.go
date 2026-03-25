// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogIdempotency is the golang structure of table catalog_idempotency for DAO operations like Where/Data.
type CatalogIdempotency struct {
	g.Meta         `orm:"table:catalog_idempotency, do:true"`
	Id             any         //
	ActorUserId    any         //
	ActionCode     any         //
	IdempotencyKey any         //
	RequestHash    any         //
	ShopNo         any         //
	SpuNo          any         //
	Status         any         //
	ResponseCode   any         //
	ResponseJson   any         //
	ExpiredAt      *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
