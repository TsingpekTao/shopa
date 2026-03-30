// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentReconciliationRecord is the golang structure for table payment_reconciliation_record.
type PaymentReconciliationRecord struct {
	Id            uint64      `json:"id"            orm:"id"             ` //
	DiffNo        string      `json:"diffNo"        orm:"diff_no"        ` //
	ReconTaskNo   string      `json:"reconTaskNo"   orm:"recon_task_no"  ` //
	DiffType      uint        `json:"diffType"      orm:"diff_type"      ` //
	PaymentNo     string      `json:"paymentNo"     orm:"payment_no"     ` //
	OrderNo       string      `json:"orderNo"       orm:"order_no"       ` //
	LocalAmount   uint64      `json:"localAmount"   orm:"local_amount"   ` //
	GatewayAmount uint64      `json:"gatewayAmount" orm:"gateway_amount" ` //
	Status        string      `json:"status"        orm:"status"         ` //
	DetailJson    string      `json:"detailJson"    orm:"detail_json"    ` //
	ResolvedAt    *gtime.Time `json:"resolvedAt"    orm:"resolved_at"    ` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` //
}
