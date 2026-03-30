// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentIdempotency is the golang structure of table payment_idempotency for DAO operations like Where/Data.
type PaymentIdempotency struct {
	g.Meta       `orm:"table:payment_idempotency, do:true"`
	Id           any         //
	IdemKey      any         //
	BizCode      any         //
	BizNo        any         //
	ResponseJson any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
