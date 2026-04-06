// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsExpireBucket is the golang structure for table points_expire_bucket.
type PointsExpireBucket struct {
	BucketNo        string      `json:"bucketNo"        orm:"bucket_no"        ` //
	UserId          uint64      `json:"userId"          orm:"user_id"          ` //
	SourceType      string      `json:"sourceType"      orm:"source_type"      ` // GRANT/RETURN/REFUND_GRACE/ADJUST
	SourceNo        string      `json:"sourceNo"        orm:"source_no"        ` //
	SourceVersion   uint64      `json:"sourceVersion"   orm:"source_version"   ` //
	BucketPeriod    string      `json:"bucketPeriod"    orm:"bucket_period"    ` // YYYYMM or custom window code
	TotalPoints     uint64      `json:"totalPoints"     orm:"total_points"     ` //
	RemainingPoints uint64      `json:"remainingPoints" orm:"remaining_points" ` //
	LockedPoints    uint64      `json:"lockedPoints"    orm:"locked_points"    ` // Currently locked points reserved by active reservations
	UsedPoints      uint64      `json:"usedPoints"      orm:"used_points"      ` //
	ExpiredPoints   uint64      `json:"expiredPoints"   orm:"expired_points"   ` //
	ReturnedPoints  uint64      `json:"returnedPoints"  orm:"returned_points"  ` //
	ExpireAt        *gtime.Time `json:"expireAt"        orm:"expire_at"        ` //
	BucketStatus    string      `json:"bucketStatus"    orm:"bucket_status"    ` // ACTIVE/EXPIRED/CLOSED
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` //
}
