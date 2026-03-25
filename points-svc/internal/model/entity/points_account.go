// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsAccount is the golang structure for table points_account.
type PointsAccount struct {
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"User ID"`                // User ID
	Balance   int64       `json:"balance"   orm:"balance"    description:"Current points balance"` // Current points balance
	Status    uint        `json:"status"    orm:"status"     description:"1 active,2 disabled"`    // 1 active,2 disabled
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`                       //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`                       //
}
