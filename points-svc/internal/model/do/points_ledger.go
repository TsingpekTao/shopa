// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsLedger is the golang structure of table points_ledger for DAO operations like Where/Data.
type PointsLedger struct {
	g.Meta       `orm:"table:points_ledger, do:true"`
	Id           any         //
	UserId       any         //
	BizType      any         // REGISTER_INIT/ORDER_PAY/REFUND/etc
	BizId        any         // Business id for idempotency
	Delta        any         // Points delta
	BalanceAfter any         // Balance snapshot after apply
	CreatedAt    *gtime.Time //
}
