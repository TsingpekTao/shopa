// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaAsset is the golang structure for table media_asset.
type MediaAsset struct {
	AssetId             uint64      `json:"assetId"             orm:"asset_id"              description:""`                                                  //
	ParentAssetId       uint64      `json:"parentAssetId"       orm:"parent_asset_id"       description:"Parent asset id, null for root ORIGINAL asset"`     // Parent asset id, null for root ORIGINAL asset
	RootAssetId         uint64      `json:"rootAssetId"         orm:"root_asset_id"         description:"Root asset id of the asset tree"`                   // Root asset id of the asset tree
	SceneCode           string      `json:"sceneCode"           orm:"scene_code"            description:""`                                                  //
	AclType             uint        `json:"aclType"             orm:"acl_type"              description:"1 PUBLIC_READ, 2 PRIVATE"`                          // 1 PUBLIC_READ, 2 PRIVATE
	AssetRole           uint        `json:"assetRole"           orm:"asset_role"            description:"1 ORIGINAL, 2 DERIVED"`                             // 1 ORIGINAL, 2 DERIVED
	DerivedKind         uint        `json:"derivedKind"         orm:"derived_kind"          description:"Derived asset subtype enum"`                        // Derived asset subtype enum
	FileName            string      `json:"fileName"            orm:"file_name"             description:""`                                                  //
	MimeType            string      `json:"mimeType"            orm:"mime_type"             description:""`                                                  //
	SizeBytes           uint64      `json:"sizeBytes"           orm:"size_bytes"            description:""`                                                  //
	ChecksumSha256      string      `json:"checksumSha256"      orm:"checksum_sha256"       description:""`                                                  //
	Etag                string      `json:"etag"                orm:"etag"                  description:""`                                                  //
	StorageProvider     string      `json:"storageProvider"     orm:"storage_provider"      description:"Internal storage provider, e.g. minio/oss"`         // Internal storage provider, e.g. minio/oss
	Bucket              string      `json:"bucket"              orm:"bucket"                description:"Internal storage bucket"`                           // Internal storage bucket
	ObjectKey           string      `json:"objectKey"           orm:"object_key"            description:"Internal object key/path"`                          // Internal object key/path
	StorageObjectHash   string      `json:"storageObjectHash"   orm:"storage_object_hash"   description:"SHA256 hash of storage_provider+bucket+object_key"` // SHA256 hash of storage_provider+bucket+object_key
	PublicUrl           string      `json:"publicUrl"           orm:"public_url"            description:"Static CDN URL for PUBLIC_READ scenes"`             // Static CDN URL for PUBLIC_READ scenes
	RiskStatus          uint        `json:"riskStatus"          orm:"risk_status"           description:""`                                                  //
	ProcessStatus       uint        `json:"processStatus"       orm:"process_status"        description:""`                                                  //
	ProcessProgress     uint        `json:"processProgress"     orm:"process_progress"      description:""`                                                  //
	ProcessErrorCode    string      `json:"processErrorCode"    orm:"process_error_code"    description:""`                                                  //
	ProcessErrorMessage string      `json:"processErrorMessage" orm:"process_error_message" description:""`                                                  //
	ProcessedAt         *gtime.Time `json:"processedAt"         orm:"processed_at"          description:""`                                                  //
	UploaderUserId      uint64      `json:"uploaderUserId"      orm:"uploader_user_id"      description:"Uploader user id from metadata"`                    // Uploader user id from metadata
	TraceBizType        string      `json:"traceBizType"        orm:"trace_biz_type"        description:"Optional trace-only biz type from InitUpload"`      // Optional trace-only biz type from InitUpload
	TraceBizNo          string      `json:"traceBizNo"          orm:"trace_biz_no"          description:"Optional trace-only biz no from InitUpload"`        // Optional trace-only biz no from InitUpload
	DeletedAt           *gtime.Time `json:"deletedAt"           orm:"deleted_at"            description:"Soft delete mark, physical delete by GC"`           // Soft delete mark, physical delete by GC
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            description:""`                                                  //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            description:""`                                                  //
}
