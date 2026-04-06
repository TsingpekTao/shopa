// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderIdempotency is the golang structure for table order_idempotency.
type OrderIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              ` //
	UserId         uint64      `json:"userId"         orm:"user_id"         ` //
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" ` //
	ActionCode     string      `json:"actionCode"     orm:"action_code"     ` //
	OrderNo        string      `json:"orderNo"        orm:"order_no"        ` //
	Status         uint        `json:"status"         orm:"status"          ` //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   ` //
	ErrorCode      string      `json:"errorCode"      orm:"error_code"      ` //
	ExpireAt       *gtime.Time `json:"expireAt"       orm:"expire_at"       ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
}
