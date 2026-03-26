// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaIdempotency 是 media_idempotency 表的结构体。
type MediaIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`                                    //
	ActorUserId    uint64      `json:"actorUserId"    orm:"actor_user_id"   description:"User id from metadata"`               // metadata 中的用户 ID。
	ActionCode     string      `json:"actionCode"     orm:"action_code"     description:"InitUpload/BatchBindAssetsToBiz/..."` // 示例：InitUpload/BatchBindAssetsToBiz/...
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" description:"x-idempotency-key"`                   // x-idempotency-key 请求头。
	RequestHash    string      `json:"requestHash"    orm:"request_hash"    description:"Request payload hash"`                // 请求有效载荷哈希。
	SceneCode      string      `json:"sceneCode"      orm:"scene_code"      description:""`                                    //
	BizType        string      `json:"bizType"        orm:"biz_type"        description:""`                                    //
	BizNo          string      `json:"bizNo"          orm:"biz_no"          description:""`                                    //
	Status         uint        `json:"status"         orm:"status"          description:""`                                    //
	ResponseCode   int         `json:"responseCode"   orm:"response_code"   description:""`                                    //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   description:""`                                    //
	ExpiredAt      *gtime.Time `json:"expiredAt"      orm:"expired_at"      description:""`                                    //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`                                    //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""`                                    //
}
