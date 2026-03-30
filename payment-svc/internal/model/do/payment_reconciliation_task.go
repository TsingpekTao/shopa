// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentReconciliationTask is the golang structure of table payment_reconciliation_task for DAO operations like Where/Data.
type PaymentReconciliationTask struct {
	g.Meta       `orm:"table:payment_reconciliation_task, do:true"`
	Id           any         //
	ReconTaskNo  any         //
	ReconDate    *gtime.Time //
	PayChannel   any         //
	Status       any         //
	TotalRecords any         //
	DiffRecords  any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
