// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewIdempotency is the golang structure for table review_idempotency.
type ReviewIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              description:""` //
	UserId         uint64      `json:"userId"         orm:"user_id"         description:""` //
	Action         string      `json:"action"         orm:"action"          description:""` //
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" description:""` //
	ResourceNo     string      `json:"resourceNo"     orm:"resource_no"     description:""` //
	Status         uint        `json:"status"         orm:"status"          description:""` //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   description:""` //
	ErrorCode      string      `json:"errorCode"      orm:"error_code"      description:""` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""` //
}
