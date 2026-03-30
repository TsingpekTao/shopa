// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SearchRebuildJob is the golang structure for table search_rebuild_job.
type SearchRebuildJob struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	JobNo        string      `json:"jobNo"        orm:"job_no"        ` //
	ReasonCode   string      `json:"reasonCode"   orm:"reason_code"   ` //
	Status       string      `json:"status"       orm:"status"        ` //
	TotalCount   uint64      `json:"totalCount"   orm:"total_count"   ` //
	SuccessCount uint64      `json:"successCount" orm:"success_count" ` //
	FailCount    uint64      `json:"failCount"    orm:"fail_count"    ` //
	StartedAt    *gtime.Time `json:"startedAt"    orm:"started_at"    ` //
	FinishedAt   *gtime.Time `json:"finishedAt"   orm:"finished_at"   ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
