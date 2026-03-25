// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaIdempotency 是 media_idempotency 表的 Go 结构体，供 DAO 的 Where/Data 等操作使用。
type MediaIdempotency struct {
	g.Meta         `orm:"table:media_idempotency, do:true"`
	Id             any         //
	ActorUserId    any         // metadata 中的用户 ID。
	ActionCode     any         // 示例：InitUpload/BatchBindAssetsToBiz/...。
	IdempotencyKey any         // x-idempotency-key 请求头。
	RequestHash    any         // 请求有效载荷哈希。
	SceneCode      any         //
	BizType        any         //
	BizNo          any         //
	Status         any         //
	ResponseCode   any         //
	ResponseJson   any         //
	ExpiredAt      *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
