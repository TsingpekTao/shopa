// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderIdempotency is the golang structure for table order_idempotency.
type OrderIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              description:""` //
	UserId         uint64      `json:"userId"         orm:"user_id"         description:""` //
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" description:""` //
	ActionCode     string      `json:"actionCode"     orm:"action_code"     description:""` //
	OrderNo        string      `json:"orderNo"        orm:"order_no"        description:""` //
	Status         uint        `json:"status"         orm:"status"          description:""` //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   description:""` //
	ErrorCode      string      `json:"errorCode"      orm:"error_code"      description:""` //
	ExpireAt       *gtime.Time `json:"expireAt"       orm:"expire_at"       description:""` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""` //
}
