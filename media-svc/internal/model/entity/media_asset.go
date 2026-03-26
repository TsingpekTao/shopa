// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaAsset 是 media_asset 表的结构体。
type MediaAsset struct {
	AssetId             uint64      `json:"assetId"             orm:"asset_id"              description:""`                                                  //
	ParentAssetId       uint64      `json:"parentAssetId"       orm:"parent_asset_id"       description:"Parent asset id, null for root ORIGINAL asset"`     // 父资产 ID，ORIGINAL 根节点时为空。
	RootAssetId         uint64      `json:"rootAssetId"         orm:"root_asset_id"         description:"Root asset id of the asset tree"`                   // 资产树的根资产 ID。
	SceneCode           string      `json:"sceneCode"           orm:"scene_code"            description:""`                                                  //
	AclType             uint        `json:"aclType"             orm:"acl_type"              description:"1 PUBLIC_READ, 2 PRIVATE"`                          // 1 表示 PUBLIC_READ，2 表示 PRIVATE。
	AssetRole           uint        `json:"assetRole"           orm:"asset_role"            description:"1 ORIGINAL, 2 DERIVED"`                             // 1 表示 ORIGINAL，2 表示 DERIVED。
	DerivedKind         uint        `json:"derivedKind"         orm:"derived_kind"          description:"Derived asset subtype enum"`                        // 派生资产子类型枚举。
	FileName            string      `json:"fileName"            orm:"file_name"             description:""`                                                  //
	MimeType            string      `json:"mimeType"            orm:"mime_type"             description:""`                                                  //
	SizeBytes           uint64      `json:"sizeBytes"           orm:"size_bytes"            description:""`                                                  //
	ChecksumSha256      string      `json:"checksumSha256"      orm:"checksum_sha256"       description:""`                                                  //
	Etag                string      `json:"etag"                orm:"etag"                  description:""`                                                  //
	StorageProvider     string      `json:"storageProvider"     orm:"storage_provider"      description:"Internal storage provider, e.g. minio/oss"`         // 内部存储提供方，例如 minio/oss。
	Bucket              string      `json:"bucket"              orm:"bucket"                description:"Internal storage bucket"`                           // 内部存储桶名。
	ObjectKey           string      `json:"objectKey"           orm:"object_key"            description:"Internal object key/path"`                          // 内部对象键/路径。
	StorageObjectHash   string      `json:"storageObjectHash"   orm:"storage_object_hash"   description:"SHA256 hash of storage_provider+bucket+object_key"` // storage_provider+bucket+object_key 的 SHA256 哈希。
	PublicUrl           string      `json:"publicUrl"           orm:"public_url"            description:"Static CDN URL for PUBLIC_READ scenes"`             // PUBLIC_READ 场景的静态 CDN URL。
	RiskStatus          uint        `json:"riskStatus"          orm:"risk_status"           description:""`                                                  //
	ProcessStatus       uint        `json:"processStatus"       orm:"process_status"        description:""`                                                  //
	ProcessProgress     uint        `json:"processProgress"     orm:"process_progress"      description:""`                                                  //
	ProcessErrorCode    string      `json:"processErrorCode"    orm:"process_error_code"    description:""`                                                  //
	ProcessErrorMessage string      `json:"processErrorMessage" orm:"process_error_message" description:""`                                                  //
	ProcessedAt         *gtime.Time `json:"processedAt"         orm:"processed_at"          description:""`                                                  //
	UploaderUserId      uint64      `json:"uploaderUserId"      orm:"uploader_user_id"      description:"Uploader user id from metadata"`                    // metadata 中的上传者 ID。
	TraceBizType        string      `json:"traceBizType"        orm:"trace_biz_type"        description:"Optional trace-only biz type from InitUpload"`      // InitUpload 传入的可选链路追踪业务类型。
	TraceBizNo          string      `json:"traceBizNo"          orm:"trace_biz_no"          description:"Optional trace-only biz no from InitUpload"`        // InitUpload 传入的可选链路追踪业务编号。
	DeletedAt           *gtime.Time `json:"deletedAt"           orm:"deleted_at"            description:"Soft delete mark, physical delete by GC"`           // 软删除标记，由 GC 负责物理删除。
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            description:""`                                                  //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            description:""`                                                  //
}
