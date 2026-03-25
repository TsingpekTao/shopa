// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaIdempotency is the golang structure of table media_idempotency for DAO operations like Where/Data.
type MediaIdempotency struct {
	g.Meta         `orm:"table:media_idempotency, do:true"`
	Id             any         //
	ActorUserId    any         // User id from metadata
	ActionCode     any         // InitUpload/BatchBindAssetsToBiz/...
	IdempotencyKey any         // x-idempotency-key
	RequestHash    any         // Request payload hash
	SceneCode      any         //
	BizType        any         //
	BizNo          any         //
	Status         any         //
	ResponseCode   any         //
	ResponseJson   any         //
	ExpiredAt      *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
