// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RefundTask is the golang structure for table refund_task.
type RefundTask struct {
	Id               uint64      `json:"id"               orm:"id"                 ` //
	RefundTaskNo     string      `json:"refundTaskNo"     orm:"refund_task_no"     ` //
	AfterSaleNo      string      `json:"afterSaleNo"      orm:"after_sale_no"      ` //
	OrderNo          string      `json:"orderNo"          orm:"order_no"           ` //
	PaymentNo        string      `json:"paymentNo"        orm:"payment_no"         ` //
	RefundAmount     uint64      `json:"refundAmount"     orm:"refund_amount"      ` //
	Status           uint        `json:"status"           orm:"status"             ` //
	RetryCount       uint        `json:"retryCount"       orm:"retry_count"        ` //
	NextRetryAt      *gtime.Time `json:"nextRetryAt"      orm:"next_retry_at"      ` //
	LastErrorCode    string      `json:"lastErrorCode"    orm:"last_error_code"    ` //
	LastErrorMessage string      `json:"lastErrorMessage" orm:"last_error_message" ` //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         ` //
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         ` //
}
