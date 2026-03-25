// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaAsset is the golang structure of table media_asset for DAO operations like Where/Data.
type MediaAsset struct {
	g.Meta              `orm:"table:media_asset, do:true"`
	AssetId             any         //
	ParentAssetId       any         // Parent asset id, null for root ORIGINAL asset
	RootAssetId         any         // Root asset id of the asset tree
	SceneCode           any         //
	AclType             any         // 1 PUBLIC_READ, 2 PRIVATE
	AssetRole           any         // 1 ORIGINAL, 2 DERIVED
	DerivedKind         any         // Derived asset subtype enum
	FileName            any         //
	MimeType            any         //
	SizeBytes           any         //
	ChecksumSha256      any         //
	Etag                any         //
	StorageProvider     any         // Internal storage provider, e.g. minio/oss
	Bucket              any         // Internal storage bucket
	ObjectKey           any         // Internal object key/path
	StorageObjectHash   any         // SHA256 hash of storage_provider+bucket+object_key
	PublicUrl           any         // Static CDN URL for PUBLIC_READ scenes
	RiskStatus          any         //
	ProcessStatus       any         //
	ProcessProgress     any         //
	ProcessErrorCode    any         //
	ProcessErrorMessage any         //
	ProcessedAt         *gtime.Time //
	UploaderUserId      any         // Uploader user id from metadata
	TraceBizType        any         // Optional trace-only biz type from InitUpload
	TraceBizNo          any         // Optional trace-only biz no from InitUpload
	DeletedAt           *gtime.Time // Soft delete mark, physical delete by GC
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
