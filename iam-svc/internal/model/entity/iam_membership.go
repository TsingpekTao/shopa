// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamMembership is the golang structure for table iam_membership.
type IamMembership struct {
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""` //
	LevelCode string      `json:"levelCode" orm:"level_code" description:""` //
	Points    uint64      `json:"points"    orm:"points"     description:""` //
	ExpireAt  *gtime.Time `json:"expireAt"  orm:"expire_at"  description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
}
