// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaScenePolicy is the golang structure of table media_scene_policy for DAO operations like Where/Data.
type MediaScenePolicy struct {
	g.Meta              `orm:"table:media_scene_policy, do:true"`
	Id                  any         //
	SceneCode           any         // Unique scene code, e.g. PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC
	AclType             any         // 1 PUBLIC_READ, 2 PRIVATE
	MaxCount            any         // Max assets per binding slot
	MaxSizeBytes        any         // Max file size in bytes, 0 means no limit
	AllowMimeJson       any         // Allowed mime list JSON
	AllowExtJson        any         // Allowed extension list JSON
	RetentionDays       any         // Retention days after unbound
	GcGraceHours        any         // Extra grace hours before physical delete
	RiskLevel           any         // 1 LOW, 2 MEDIUM, 3 HIGH
	RiskAsyncEnabled    any         // Whether to enable async risk scanning
	ProcessAsyncEnabled any         // Whether to enable async document processing
	Status              any         // 1 ENABLED, 2 DISABLED
	Remark              any         //
	ExtJson             any         //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
