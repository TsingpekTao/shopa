// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentReconciliationRecord is the golang structure of table payment_reconciliation_record for DAO operations like Where/Data.
type PaymentReconciliationRecord struct {
	g.Meta        `orm:"table:payment_reconciliation_record, do:true"`
	Id            any         //
	DiffNo        any         //
	ReconTaskNo   any         //
	DiffType      any         //
	PaymentNo     any         //
	OrderNo       any         //
	LocalAmount   any         //
	GatewayAmount any         //
	Status        any         //
	DetailJson    any         //
	ResolvedAt    *gtime.Time //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
