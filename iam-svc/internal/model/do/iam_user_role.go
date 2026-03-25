// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserRole is the golang structure of table iam_user_role for DAO operations like Where/Data.
type IamUserRole struct {
	g.Meta    `orm:"table:iam_user_role, do:true"`
	Id        any         //
	UserId    any         //
	RoleCode  any         // 1 customer,2 seller,3 admin,4 cs
	ScopeType any         // 1 global,2 shop
	ScopeId   any         //
	Status    any         // 1 active,2 disabled
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
