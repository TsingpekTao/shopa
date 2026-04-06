// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsAccount is the golang structure for table points_account.
type PointsAccount struct {
	UserId              uint64      `json:"userId"              orm:"user_id"               ` // User ID
	AvailableBalance    int64       `json:"availableBalance"    orm:"available_balance"     ` // Current available points balance, may be negative when debt exists
	FrozenBalance       int64       `json:"frozenBalance"       orm:"frozen_balance"        ` // Currently frozen points balance
	StatusCode          string      `json:"statusCode"          orm:"status_code"           ` // ACTIVE/FROZEN/DISABLED
	TotalEarnedPoints   uint64      `json:"totalEarnedPoints"   orm:"total_earned_points"   ` // Lifetime granted points
	TotalUsedPoints     uint64      `json:"totalUsedPoints"     orm:"total_used_points"     ` // Lifetime confirmed spent points
	TotalExpiredPoints  uint64      `json:"totalExpiredPoints"  orm:"total_expired_points"  ` // Lifetime expired points
	TotalAdjustedPoints int64       `json:"totalAdjustedPoints" orm:"total_adjusted_points" ` // Lifetime manual adjustment points
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            ` //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            ` //
}
