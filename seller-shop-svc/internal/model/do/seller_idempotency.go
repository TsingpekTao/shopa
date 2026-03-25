// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerIdempotency is the golang structure of table seller_idempotency for DAO operations like Where/Data.
type SellerIdempotency struct {
	g.Meta         `orm:"table:seller_idempotency, do:true"`
	Id             any         //
	ActorUserId    any         // User id or operator id
	ActionCode     any         // CreateApplicationDraft/ApproveApplication/...
	IdempotencyKey any         // x-idempotency-key
	RequestHash    any         // Payload hash for mismatch detection
	ResourceType   any         // application/shop/entity
	ResourceNo     any         // Business id
	Status         any         // 1 PROCESSING,2 SUCCEEDED,3 FAILED
	ResponseCode   any         //
	ResponseJson   any         //
	ExpiredAt      *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
