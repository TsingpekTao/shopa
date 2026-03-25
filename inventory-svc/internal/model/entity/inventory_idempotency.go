// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryIdempotency is the golang structure for table inventory_idempotency.
type InventoryIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              ` //
	ActorScope     string      `json:"actorScope"     orm:"actor_scope"     ` // SELLER/ORDER/ADMIN/SYSTEM
	ActorId        uint64      `json:"actorId"        orm:"actor_id"        ` //
	ActionCode     string      `json:"actionCode"     orm:"action_code"     ` //
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" ` //
	RequestHash    string      `json:"requestHash"    orm:"request_hash"    ` //
	Status         uint        `json:"status"         orm:"status"          ` //
	ResponseCode   int         `json:"responseCode"   orm:"response_code"   ` //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   ` //
	ExpiredAt      *gtime.Time `json:"expiredAt"      orm:"expired_at"      ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
}
