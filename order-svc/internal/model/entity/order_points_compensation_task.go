// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderPointsCompensationTask is the golang structure for table order_points_compensation_task.
type OrderPointsCompensationTask struct {
	Id                  uint64      `json:"id"                  orm:"id"                    ` //
	TaskNo              string      `json:"taskNo"              orm:"task_no"               ` //
	OrderNo             string      `json:"orderNo"             orm:"order_no"              ` //
	UserId              uint64      `json:"userId"              orm:"user_id"               ` //
	PointsReservationNo string      `json:"pointsReservationNo" orm:"points_reservation_no" ` //
	ActionCode          string      `json:"actionCode"          orm:"action_code"           ` //
	TaskStatus          string      `json:"taskStatus"          orm:"task_status"           ` //
	RetryCount          uint        `json:"retryCount"          orm:"retry_count"           ` //
	NextRetryAt         *gtime.Time `json:"nextRetryAt"         orm:"next_retry_at"         ` //
	LastError           string      `json:"lastError"           orm:"last_error"            ` //
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            ` //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            ` //
}
