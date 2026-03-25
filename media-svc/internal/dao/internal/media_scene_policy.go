// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaScenePolicyDao is the data access object for the table media_scene_policy.
type MediaScenePolicyDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  MediaScenePolicyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// MediaScenePolicyColumns defines and stores column names for the table media_scene_policy.
type MediaScenePolicyColumns struct {
	Id                  string //
	SceneCode           string // Unique scene code, e.g. PRODUCT_MAIN/SELLER_CERT/AGENT_CHAT_DOC
	AclType             string // 1 PUBLIC_READ, 2 PRIVATE
	MaxCount            string // Max assets per binding slot
	MaxSizeBytes        string // Max file size in bytes, 0 means no limit
	AllowMimeJson       string // Allowed mime list JSON
	AllowExtJson        string // Allowed extension list JSON
	RetentionDays       string // Retention days after unbound
	GcGraceHours        string // Extra grace hours before physical delete
	RiskLevel           string // 1 LOW, 2 MEDIUM, 3 HIGH
	RiskAsyncEnabled    string // Whether to enable async risk scanning
	ProcessAsyncEnabled string // Whether to enable async document processing
	Status              string // 1 ENABLED, 2 DISABLED
	Remark              string //
	ExtJson             string //
	CreatedAt           string //
	UpdatedAt           string //
}

// mediaScenePolicyColumns holds the columns for the table media_scene_policy.
var mediaScenePolicyColumns = MediaScenePolicyColumns{
	Id:                  "id",
	SceneCode:           "scene_code",
	AclType:             "acl_type",
	MaxCount:            "max_count",
	MaxSizeBytes:        "max_size_bytes",
	AllowMimeJson:       "allow_mime_json",
	AllowExtJson:        "allow_ext_json",
	RetentionDays:       "retention_days",
	GcGraceHours:        "gc_grace_hours",
	RiskLevel:           "risk_level",
	RiskAsyncEnabled:    "risk_async_enabled",
	ProcessAsyncEnabled: "process_async_enabled",
	Status:              "status",
	Remark:              "remark",
	ExtJson:             "ext_json",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewMediaScenePolicyDao creates and returns a new DAO object for table data access.
func NewMediaScenePolicyDao(handlers ...gdb.ModelHandler) *MediaScenePolicyDao {
	return &MediaScenePolicyDao{
		group:    "default",
		table:    "media_scene_policy",
		columns:  mediaScenePolicyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MediaScenePolicyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MediaScenePolicyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MediaScenePolicyDao) Columns() MediaScenePolicyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MediaScenePolicyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MediaScenePolicyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MediaScenePolicyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
