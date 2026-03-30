// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SearchRebuildJob is the golang structure of table search_rebuild_job for DAO operations like Where/Data.
type SearchRebuildJob struct {
	g.Meta       `orm:"table:search_rebuild_job, do:true"`
	Id           any         //
	JobNo        any         //
	ReasonCode   any         //
	Status       any         //
	TotalCount   any         //
	SuccessCount any         //
	FailCount    any         //
	StartedAt    *gtime.Time //
	FinishedAt   *gtime.Time //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
