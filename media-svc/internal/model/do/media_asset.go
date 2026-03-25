// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaAsset 是 media_asset 表的 Go 结构体，供 DAO 的 Where/Data 等操作使用。
type MediaAsset struct {
	g.Meta              `orm:"table:media_asset, do:true"`
	AssetId             any         //
	ParentAssetId       any         // 父资产 ID，ORIGINAL 根节点时为 0。
	RootAssetId         any         // 资产树的根资产 ID。
	SceneCode           any         //
	AclType             any         // 1 表示 PUBLIC_READ，2 表示 PRIVATE。
	AssetRole           any         // 1 表示 ORIGINAL，2 表示 DERIVED。
	DerivedKind         any         // 派生资产子类型枚举。
	FileName            any         //
	MimeType            any         //
	SizeBytes           any         //
	ChecksumSha256      any         //
	Etag                any         //
	StorageProvider     any         // 内部存储提供方，例如 minio/oss。
	Bucket              any         // 内部存储桶名。
	ObjectKey           any         // 内部对象键/路径。
	StorageObjectHash   any         // storage_provider+bucket+object_key 的 SHA256 哈希。
	PublicUrl           any         // PUBLIC_READ 场景的静态 CDN 访问地址。
	RiskStatus          any         //
	ProcessStatus       any         //
	ProcessProgress     any         //
	ProcessErrorCode    any         //
	ProcessErrorMessage any         //
	ProcessedAt         *gtime.Time //
	UploaderUserId      any         // metadata 中的上传者 ID。
	TraceBizType        any         // InitUpload 传入的可选链路追踪业务类型。
	TraceBizNo          any         // InitUpload 传入的可选链路追踪业务编号。
	DeletedAt           *gtime.Time // 软删除标记，由 GC 负责物理删除。
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
