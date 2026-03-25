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
	g.Meta    `orm:"table:points_account, do:true"`
	UserId    any         // User ID
	Balance   any         // Current points balance
	Status    any         // 1 active,2 disabled
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
