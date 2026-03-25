// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaIdempotencyDao is the data access object for the table media_idempotency.
type MediaIdempotencyDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  MediaIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// MediaIdempotencyColumns defines and stores column names for the table media_idempotency.
type MediaIdempotencyColumns struct {
	Id             string //
	ActorUserId    string // User id from metadata
	ActionCode     string // InitUpload/BatchBindAssetsToBiz/...
	IdempotencyKey string // x-idempotency-key
	RequestHash    string // Request payload hash
	SceneCode      string //
	BizType        string //
	BizNo          string //
	Status         string //
	ResponseCode   string //
	ResponseJson   string //
	ExpiredAt      string //
	CreatedAt      string //
	UpdatedAt      string //
}

// mediaIdempotencyColumns holds the columns for the table media_idempotency.
var mediaIdempotencyColumns = MediaIdempotencyColumns{
	Id:             "id",
	ActorUserId:    "actor_user_id",
	ActionCode:     "action_code",
	IdempotencyKey: "idempotency_key",
	RequestHash:    "request_hash",
	SceneCode:      "scene_code",
	BizType:        "biz_type",
	BizNo:          "biz_no",
	Status:         "status",
	ResponseCode:   "response_code",
	ResponseJson:   "response_json",
	ExpiredAt:      "expired_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewMediaIdempotencyDao creates and returns a new DAO object for table data access.
func NewMediaIdempotencyDao(handlers ...gdb.ModelHandler) *MediaIdempotencyDao {
	return &MediaIdempotencyDao{
		group:    "default",
		table:    "media_idempotency",
		columns:  mediaIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MediaIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MediaIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MediaIdempotencyDao) Columns() MediaIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MediaIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MediaIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MediaIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
