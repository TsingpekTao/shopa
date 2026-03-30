// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentReconciliationTask is the golang structure for table payment_reconciliation_task.
type PaymentReconciliationTask struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	ReconTaskNo  string      `json:"reconTaskNo"  orm:"recon_task_no" ` //
	ReconDate    *gtime.Time `json:"reconDate"    orm:"recon_date"    ` //
	PayChannel   uint        `json:"payChannel"   orm:"pay_channel"   ` //
	Status       string      `json:"status"       orm:"status"        ` //
	TotalRecords uint64      `json:"totalRecords" orm:"total_records" ` //
	DiffRecords  uint64      `json:"diffRecords"  orm:"diff_records"  ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
