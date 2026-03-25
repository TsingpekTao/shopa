// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogIdempotency is the golang structure for table catalog_idempotency.
type CatalogIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              ` //
	ActorUserId    uint64      `json:"actorUserId"    orm:"actor_user_id"   ` //
	ActionCode     string      `json:"actionCode"     orm:"action_code"     ` //
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" ` //
	RequestHash    string      `json:"requestHash"    orm:"request_hash"    ` //
	ShopNo         string      `json:"shopNo"         orm:"shop_no"         ` //
	SpuNo          string      `json:"spuNo"          orm:"spu_no"          ` //
	Status         uint        `json:"status"         orm:"status"          ` //
	ResponseCode   int         `json:"responseCode"   orm:"response_code"   ` //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   ` //
	ExpiredAt      *gtime.Time `json:"expiredAt"      orm:"expired_at"      ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
}
