// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamMembership is the golang structure of table iam_membership for DAO operations like Where/Data.
type IamMembership struct {
	g.Meta    `orm:"table:iam_membership, do:true"`
	UserId    any         //
	LevelCode any         //
	Points    any         //
	ExpireAt  *gtime.Time //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
