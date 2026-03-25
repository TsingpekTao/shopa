// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaAssetDao is the data access object for the table media_asset.
type MediaAssetDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MediaAssetColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MediaAssetColumns defines and stores column names for the table media_asset.
type MediaAssetColumns struct {
	AssetId             string //
	ParentAssetId       string // Parent asset id, null for root ORIGINAL asset
	RootAssetId         string // Root asset id of the asset tree
	SceneCode           string //
	AclType             string // 1 PUBLIC_READ, 2 PRIVATE
	AssetRole           string // 1 ORIGINAL, 2 DERIVED
	DerivedKind         string // Derived asset subtype enum
	FileName            string //
	MimeType            string //
	SizeBytes           string //
	ChecksumSha256      string //
	Etag                string //
	StorageProvider     string // Internal storage provider, e.g. minio/oss
	Bucket              string // Internal storage bucket
	ObjectKey           string // Internal object key/path
	StorageObjectHash   string // SHA256 hash of storage_provider+bucket+object_key
	PublicUrl           string // Static CDN URL for PUBLIC_READ scenes
	RiskStatus          string //
	ProcessStatus       string //
	ProcessProgress     string //
	ProcessErrorCode    string //
	ProcessErrorMessage string //
	ProcessedAt         string //
	UploaderUserId      string // Uploader user id from metadata
	TraceBizType        string // Optional trace-only biz type from InitUpload
	TraceBizNo          string // Optional trace-only biz no from InitUpload
	DeletedAt           string // Soft delete mark, physical delete by GC
	CreatedAt           string //
	UpdatedAt           string //
}

// mediaAssetColumns holds the columns for the table media_asset.
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

// NewMediaAssetDao creates and returns a new DAO object for table data access.
func NewMediaAssetDao(handlers ...gdb.ModelHandler) *MediaAssetDao {
	return &MediaAssetDao{
		group:    "default",
		table:    "media_asset",
		columns:  mediaAssetColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MediaAssetDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MediaAssetDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MediaAssetDao) Columns() MediaAssetColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MediaAssetDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MediaAssetDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *MediaAssetDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
