// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsLedger is the golang structure for table points_ledger.
type PointsLedger struct {
	Id           uint64      `json:"id"           orm:"id"            description:""`                                   //
	UserId       uint64      `json:"userId"       orm:"user_id"       description:""`                                   //
	BizType      string      `json:"bizType"      orm:"biz_type"      description:"REGISTER_INIT/ORDER_PAY/REFUND/etc"` // REGISTER_INIT/ORDER_PAY/REFUND/etc
	BizId        string      `json:"bizId"        orm:"biz_id"        description:"Business id for idempotency"`        // Business id for idempotency
	Delta        int64       `json:"delta"        orm:"delta"         description:"Points delta"`                       // Points delta
	BalanceAfter int64       `json:"balanceAfter" orm:"balance_after" description:"Balance snapshot after apply"`       // Balance snapshot after apply
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`                                   //
}
