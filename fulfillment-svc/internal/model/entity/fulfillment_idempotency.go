// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentIdempotency is the golang structure for table fulfillment_idempotency.
type FulfillmentIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`                                  //
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" description:""`                                  //
	Scope          string      `json:"scope"          orm:"scope"           description:""`                                  //
	RequestHash    string      `json:"requestHash"    orm:"request_hash"    description:""`                                  //
	ResourceId     string      `json:"resourceId"     orm:"resource_id"     description:""`                                  //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   description:""`                                  //
	Status         uint        `json:"status"         orm:"status"          description:"1 SUCCEEDED 2 PROCESSING 3 FAILED"` // 1 SUCCEEDED 2 PROCESSING 3 FAILED
	ExpireAt       *gtime.Time `json:"expireAt"       orm:"expire_at"       description:""`                                  //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`                                  //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""`                                  //
}
