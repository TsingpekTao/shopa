// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerIdempotency is the golang structure for table seller_idempotency.
type SellerIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              ` //
	ActorUserId    uint64      `json:"actorUserId"    orm:"actor_user_id"   ` // User id or operator id
	ActionCode     string      `json:"actionCode"     orm:"action_code"     ` // CreateApplicationDraft/ApproveApplication/...
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" ` // x-idempotency-key
	RequestHash    string      `json:"requestHash"    orm:"request_hash"    ` // Payload hash for mismatch detection
	ResourceType   string      `json:"resourceType"   orm:"resource_type"   ` // application/shop/entity
	ResourceNo     string      `json:"resourceNo"     orm:"resource_no"     ` // Business id
	Status         uint        `json:"status"         orm:"status"          ` // 1 PROCESSING,2 SUCCEEDED,3 FAILED
	ResponseCode   int         `json:"responseCode"   orm:"response_code"   ` //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   ` //
	ExpiredAt      *gtime.Time `json:"expiredAt"      orm:"expired_at"      ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
}
