// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaIdempotency is the golang structure for table media_idempotency.
type MediaIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`                                    //
	ActorUserId    uint64      `json:"actorUserId"    orm:"actor_user_id"   description:"User id from metadata"`               // User id from metadata
	ActionCode     string      `json:"actionCode"     orm:"action_code"     description:"InitUpload/BatchBindAssetsToBiz/..."` // InitUpload/BatchBindAssetsToBiz/...
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" description:"x-idempotency-key"`                   // x-idempotency-key
	RequestHash    string      `json:"requestHash"    orm:"request_hash"    description:"Request payload hash"`                // Request payload hash
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
