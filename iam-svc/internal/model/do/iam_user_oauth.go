// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserOauth is the golang structure of table iam_user_oauth for DAO operations like Where/Data.
type IamUserOauth struct {
	g.Meta      `orm:"table:iam_user_oauth, do:true"`
	Id          any         //
	UserId      any         //
	Provider    any         // 1 wechat,2 alipay
	ProviderUid any         //
	UnionId     any         //
	Status      any         // 1 active,2 unbound
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
