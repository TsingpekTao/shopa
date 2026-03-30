// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentIdempotency is the golang structure for table payment_idempotency.
type PaymentIdempotency struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	IdemKey      string      `json:"idemKey"      orm:"idem_key"      ` //
	BizCode      string      `json:"bizCode"      orm:"biz_code"      ` //
	BizNo        string      `json:"bizNo"        orm:"biz_no"        ` //
	ResponseJson string      `json:"responseJson" orm:"response_json" ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
