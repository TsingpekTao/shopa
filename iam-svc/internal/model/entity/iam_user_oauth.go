// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IamUserOauth is the golang structure for table iam_user_oauth.
type IamUserOauth struct {
	Id          uint64      `json:"id"          orm:"id"           description:""`                   //
	UserId      uint64      `json:"userId"      orm:"user_id"      description:""`                   //
	Provider    uint        `json:"provider"    orm:"provider"     description:"1 wechat,2 alipay"`  // 1 wechat,2 alipay
	ProviderUid string      `json:"providerUid" orm:"provider_uid" description:""`                   //
	UnionId     string      `json:"unionId"     orm:"union_id"     description:""`                   //
	Status      uint        `json:"status"      orm:"status"       description:"1 active,2 unbound"` // 1 active,2 unbound
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`                   //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`                   //
}
