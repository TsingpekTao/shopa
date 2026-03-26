// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaScenePolicy 是 media_scene_policy 表的 Go 结构体，供 DAO 的 Where/Data 等操作使用。
type MediaScenePolicy struct {
	g.Meta              `orm:"table:media_scene_policy, do:true"`
	Id                  any         //
	SceneCode           any         // 唯一场景码，如 PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC。
	AclType             any         // 1 表示 PUBLIC_READ，2 表示 PRIVATE。
	MaxCount            any         // 每个绑定槽允许的最大资产数量。
	MaxSizeBytes        any         // 最大文件大小（字节），0 表示不限制。
	AllowMimeJson       any         // 允许的 MIME 列表 JSON。
	AllowExtJson        any         // 允许的扩展名列表 JSON。
	RetentionDays       any         // 解绑后保留天数。
	GcGraceHours        any         // 物理删除前的额外宽限小时数。
	RiskLevel           any         // 1 表示 LOW，2 表示 MEDIUM，3 表示 HIGH。
	RiskAsyncEnabled    any         // 是否启用异步风控扫描。
	ProcessAsyncEnabled any         // 是否启用异步处理。
	Status              any         // 1 表示 ENABLED，2 表示 DISABLED。
	Remark              any         //
	ExtJson             any         //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
