// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChatIdempotency is the golang structure for table chat_idempotency.
type ChatIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              description:""` //
	UserId         uint64      `json:"userId"         orm:"user_id"         description:""` //
	BizCode        string      `json:"bizCode"        orm:"biz_code"        description:""` //
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" description:""` //
	TargetNo       string      `json:"targetNo"       orm:"target_no"       description:""` //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   description:""` //
	Status         int         `json:"status"         orm:"status"          description:""` //
	ExpiredAt      *gtime.Time `json:"expiredAt"      orm:"expired_at"      description:""` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""` //
}
