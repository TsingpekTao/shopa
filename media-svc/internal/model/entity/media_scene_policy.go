// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaScenePolicy is the golang structure for table media_scene_policy.
type MediaScenePolicy struct {
	Id                  uint64      `json:"id"                  orm:"id"                    description:""`                                                                //
	SceneCode           string      `json:"sceneCode"           orm:"scene_code"            description:"Unique scene code, e.g. PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC"` // Unique scene code, e.g. PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC
	AclType             uint        `json:"aclType"             orm:"acl_type"              description:"1 PUBLIC_READ, 2 PRIVATE"`                                        // 1 PUBLIC_READ, 2 PRIVATE
	MaxCount            uint        `json:"maxCount"            orm:"max_count"             description:"Max assets per binding slot"`                                     // Max assets per binding slot
	MaxSizeBytes        uint64      `json:"maxSizeBytes"        orm:"max_size_bytes"        description:"Max file size in bytes, 0 means no limit"`                        // Max file size in bytes, 0 means no limit
	AllowMimeJson       string      `json:"allowMimeJson"       orm:"allow_mime_json"       description:"Allowed mime list JSON"`                                          // Allowed mime list JSON
	AllowExtJson        string      `json:"allowExtJson"        orm:"allow_ext_json"        description:"Allowed extension list JSON"`                                     // Allowed extension list JSON
	RetentionDays       uint        `json:"retentionDays"       orm:"retention_days"        description:"Retention days after unbound"`                                    // Retention days after unbound
	GcGraceHours        uint        `json:"gcGraceHours"        orm:"gc_grace_hours"        description:"Extra grace hours before physical delete"`                        // Extra grace hours before physical delete
	RiskLevel           uint        `json:"riskLevel"           orm:"risk_level"            description:"1 LOW, 2 MEDIUM, 3 HIGH"`                                         // 1 LOW, 2 MEDIUM, 3 HIGH
	RiskAsyncEnabled    int         `json:"riskAsyncEnabled"    orm:"risk_async_enabled"    description:"Whether to enable async risk scanning"`                           // Whether to enable async risk scanning
	ProcessAsyncEnabled int         `json:"processAsyncEnabled" orm:"process_async_enabled" description:"Whether to enable async document processing"`                     // Whether to enable async document processing
	Status              uint        `json:"status"              orm:"status"                description:"1 ENABLED, 2 DISABLED"`                                           // 1 ENABLED, 2 DISABLED
	Remark              string      `json:"remark"              orm:"remark"                description:""`                                                                //
	ExtJson             string      `json:"extJson"             orm:"ext_json"              description:""`                                                                //
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            description:""`                                                                //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            description:""`                                                                //
}
