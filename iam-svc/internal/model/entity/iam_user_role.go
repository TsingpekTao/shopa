// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserRole is the golang structure for table iam_user_role.
type IamUserRole struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`                                 //
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`                                 //
	RoleCode  uint        `json:"roleCode"  orm:"role_code"  description:"1 customer,2 seller,3 admin,4 cs"` // 1 customer,2 seller,3 admin,4 cs
	ScopeType uint        `json:"scopeType" orm:"scope_type" description:"1 global,2 shop"`                  // 1 global,2 shop
	ScopeId   uint64      `json:"scopeId"   orm:"scope_id"   description:""`                                 //
	Status    uint        `json:"status"    orm:"status"     description:"1 active,2 disabled"`              // 1 active,2 disabled
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`                                 //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`                                 //
}
