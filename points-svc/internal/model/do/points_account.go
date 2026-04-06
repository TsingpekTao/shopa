// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsAccount is the golang structure of table points_account for DAO operations like Where/Data.
type PointsAccount struct {
	g.Meta              `orm:"table:points_account, do:true"`
	UserId              any         // User ID
	AvailableBalance    any         // Current available points balance, may be negative when debt exists
	FrozenBalance       any         // Currently frozen points balance
	StatusCode          any         // ACTIVE/FROZEN/DISABLED
	TotalEarnedPoints   any         // Lifetime granted points
	TotalUsedPoints     any         // Lifetime confirmed spent points
	TotalExpiredPoints  any         // Lifetime expired points
	TotalAdjustedPoints any         // Lifetime manual adjustment points
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
