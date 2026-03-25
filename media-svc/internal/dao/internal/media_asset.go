// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaAssetDao 是 media_asset 表的数据访问对象。
type MediaAssetDao struct {
	table    string             // table 表示 DAO 底层对应的表名。
	group    string             // group 表示当前 DAO 的数据库配置分组。
	columns  MediaAssetColumns  // columns 包含表的全部列名，方便直接引用。
	handlers []gdb.ModelHandler // handlers 用于自定义数据模型的修改处理。
}

// MediaAssetColumns 定义并存储 media_asset 表的列名。
type MediaAssetColumns struct {
	AssetId             string //
	ParentAssetId       string // 父资产 ID，ORIGINAL 根节点时为 0。
	RootAssetId         string // 资产树的根资产 ID。
	SceneCode           string //
	AclType             string // 1 表示 PUBLIC_READ，2 表示 PRIVATE。
	AssetRole           string // 1 表示 ORIGINAL，2 表示 DERIVED。
	DerivedKind         string // 派生资产子类型枚举。
	FileName            string //
	MimeType            string //
	SizeBytes           string //
	ChecksumSha256      string //
	Etag                string //
	StorageProvider     string // 内部存储提供方，例如 minio/oss。
	Bucket              string // 内部存储桶名。
	ObjectKey           string // 内部对象键/路径。
	StorageObjectHash   string // storage_provider+bucket+object_key 的 SHA256 哈希。
	PublicUrl           string // PUBLIC_READ 场景的静态 CDN URL。
	RiskStatus          string //
	ProcessStatus       string //
	ProcessProgress     string //
	ProcessErrorCode    string //
	ProcessErrorMessage string //
	ProcessedAt         string //
	UploaderUserId      string // metadata 中的上传者 ID。
	TraceBizType        string // InitUpload 传入的可选链路追踪业务类型。
	TraceBizNo          string // InitUpload 传入的可选链路追踪业务编号。
	DeletedAt           string // 软删除标记，由 GC 物理删除。
	CreatedAt           string //
	UpdatedAt           string //
}

// mediaAssetColumns 保存 media_asset 表的列名映射。
var mediaAssetColumns = MediaAssetColumns{
	AssetId:             "asset_id",
	ParentAssetId:       "parent_asset_id",
	RootAssetId:         "root_asset_id",
	SceneCode:           "scene_code",
	AclType:             "acl_type",
	AssetRole:           "asset_role",
	DerivedKind:         "derived_kind",
	FileName:            "file_name",
	MimeType:            "mime_type",
	SizeBytes:           "size_bytes",
	ChecksumSha256:      "checksum_sha256",
	Etag:                "etag",
	StorageProvider:     "storage_provider",
	Bucket:              "bucket",
	ObjectKey:           "object_key",
	StorageObjectHash:   "storage_object_hash",
	PublicUrl:           "public_url",
	RiskStatus:          "risk_status",
	ProcessStatus:       "process_status",
	ProcessProgress:     "process_progress",
	ProcessErrorCode:    "process_error_code",
	ProcessErrorMessage: "process_error_message",
	ProcessedAt:         "processed_at",
	UploaderUserId:      "uploader_user_id",
	TraceBizType:        "trace_biz_type",
	TraceBizNo:          "trace_biz_no",
	DeletedAt:           "deleted_at",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewMediaAssetDao 创建并返回 media_asset 表的数据访问对象。
func NewMediaAssetDao(handlers ...gdb.ModelHandler) *MediaAssetDao {
	return &MediaAssetDao{
		group:    "default",
		table:    "media_asset",
		columns:  mediaAssetColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 的底层数据库管理对象。
func (dao *MediaAssetDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 操作的表名。
func (dao *MediaAssetDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *MediaAssetDao) Columns() MediaAssetColumns {
	return dao.columns
}

// Group 返回 DAO 使用的数据库配置分组名称。
func (dao *MediaAssetDao) Group() string {
	return dao.group
}

// Ctx 创建并返回当前 DAO 的 Model，并自动设置操作上下文。
func (dao *MediaAssetDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction 用参数函数 f 包裹事务逻辑。
// 如果 f 返回非 nil 错误则回滚并返回该错误。
// 如果 f 返回 nil 则提交并返回 nil。
//
// 注意：f 内请勿显式提交或回滚，事务由此方法自动管理。
func (dao *MediaAssetDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
