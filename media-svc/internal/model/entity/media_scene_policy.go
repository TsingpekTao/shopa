// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaScenePolicy 是 media_scene_policy 表的结构体。
type MediaScenePolicy struct {
	Id                  uint64      `json:"id"                  orm:"id"                    description:""`                                                                //
	SceneCode           string      `json:"sceneCode"           orm:"scene_code"            description:"Unique scene code, e.g. PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC"` // 唯一场景码，例如 PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC。
	AclType             uint        `json:"aclType"             orm:"acl_type"              description:"1 PUBLIC_READ, 2 PRIVATE"`                                        // 1 表示 PUBLIC_READ，2 表示 PRIVATE。
	MaxCount            uint        `json:"maxCount"            orm:"max_count"             description:"Max assets per binding slot"`                                     // 每个绑定槽允许的最大资产数量。
	MaxSizeBytes        uint64      `json:"maxSizeBytes"        orm:"max_size_bytes"        description:"Max file size in bytes, 0 means no limit"`                        // 最大文件大小（字节），0 表示不限制。
	AllowMimeJson       string      `json:"allowMimeJson"       orm:"allow_mime_json"       description:"Allowed mime list JSON"`                                          // 允许的 MIME 列表 JSON。
	AllowExtJson        string      `json:"allowExtJson"        orm:"allow_ext_json"        description:"Allowed extension list JSON"`                                     // 允许的扩展名列表 JSON。
	RetentionDays       uint        `json:"retentionDays"       orm:"retention_days"        description:"Retention days after unbound"`                                    // 解绑后保留天数。
	GcGraceHours        uint        `json:"gcGraceHours"        orm:"gc_grace_hours"        description:"Extra grace hours before physical delete"`                        // 物理删除前的额外宽限小时数。
	RiskLevel           uint        `json:"riskLevel"           orm:"risk_level"            description:"1 LOW, 2 MEDIUM, 3 HIGH"`                                         // 1 表示 LOW，2 表示 MEDIUM，3 表示 HIGH。
	RiskAsyncEnabled    int         `json:"riskAsyncEnabled"    orm:"risk_async_enabled"    description:"Whether to enable async risk scanning"`                           // 是否启用异步风控扫描。
	ProcessAsyncEnabled int         `json:"processAsyncEnabled" orm:"process_async_enabled" description:"Whether to enable async document processing"`                     // 是否启用异步处理。
	Status              uint        `json:"status"              orm:"status"                description:"1 ENABLED, 2 DISABLED"`                                           // 1 表示 ENABLED，2 表示 DISABLED。
	Remark              string      `json:"remark"              orm:"remark"                description:""`                                                                //
	ExtJson             string      `json:"extJson"             orm:"ext_json"              description:""`                                                                //
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            description:""`                                                                //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            description:""`                                                                //
}
