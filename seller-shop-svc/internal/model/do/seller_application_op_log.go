// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerApplicationOpLog is the golang structure of table seller_application_op_log for DAO operations like Where/Data.
type SellerApplicationOpLog struct {
	g.Meta         `orm:"table:seller_application_op_log, do:true"`
	Id             any         //
	ApplicationNo  any         //
	OwnerUserId    any         //
	OperatorUserId any         // Admin or seller operator
	ActionCode     any         // CREATE_DRAFT/SUBMIT/APPROVE/REJECT/RESUBMIT
	FromStatus     any         //
	ToStatus       any         //
	Comment        any         //
	ExtraJson      any         //
	CreatedAt      *gtime.Time //
}
