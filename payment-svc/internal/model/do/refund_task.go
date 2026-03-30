// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RefundTask is the golang structure of table refund_task for DAO operations like Where/Data.
type RefundTask struct {
	g.Meta           `orm:"table:refund_task, do:true"`
	Id               any         //
	RefundTaskNo     any         //
	AfterSaleNo      any         //
	OrderNo          any         //
	PaymentNo        any         //
	RefundAmount     any         //
	Status           any         //
	RetryCount       any         //
	NextRetryAt      *gtime.Time //
	LastErrorCode    any         //
	LastErrorMessage any         //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
}
