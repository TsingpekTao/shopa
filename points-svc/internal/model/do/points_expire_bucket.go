// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsExpireBucket is the golang structure of table points_expire_bucket for DAO operations like Where/Data.
type PointsExpireBucket struct {
	g.Meta          `orm:"table:points_expire_bucket, do:true"`
	BucketNo        any         //
	UserId          any         //
	SourceType      any         // GRANT/RETURN/REFUND_GRACE/ADJUST
	SourceNo        any         //
	SourceVersion   any         //
	BucketPeriod    any         // YYYYMM or custom window code
	TotalPoints     any         //
	RemainingPoints any         //
	LockedPoints    any         // Currently locked points reserved by active reservations
	UsedPoints      any         //
	ExpiredPoints   any         //
	ReturnedPoints  any         //
	ExpireAt        *gtime.Time //
	BucketStatus    any         // ACTIVE/EXPIRED/CLOSED
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
