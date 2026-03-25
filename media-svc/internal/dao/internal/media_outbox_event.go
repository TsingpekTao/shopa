// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaOutboxEventDao is the data access object for the table media_outbox_event.
type MediaOutboxEventDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  MediaOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// MediaOutboxEventColumns defines and stores column names for the table media_outbox_event.
type MediaOutboxEventColumns struct {
	Id            string //
	EventId       string // Global unique event id
	EventType     string // AssetProcessingCompleted/AssetProcessingFailed/...
	AggregateType string // ASSET/BINDING/SCENE
	AggregateKey  string // asset_id/biz_key/scene_code
	RequestId     string //
	PayloadJson   string //
	Status        string //
	AvailableAt   string //
	SentAt        string //
	FailCount     string //
	LastError     string //
	CreatedAt     string //
	UpdatedAt     string //
}

// mediaOutboxEventColumns holds the columns for the table media_outbox_event.
var mediaOutboxEventColumns = MediaOutboxEventColumns{
	Id:            "id",
	EventId:       "event_id",
	EventType:     "event_type",
	AggregateType: "aggregate_type",
	AggregateKey:  "aggregate_key",
	RequestId:     "request_id",
	PayloadJson:   "payload_json",
	Status:        "status",
	AvailableAt:   "available_at",
	SentAt:        "sent_at",
	FailCount:     "fail_count",
	LastError:     "last_error",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewMediaOutboxEventDao creates and returns a new DAO object for table data access.
func NewMediaOutboxEventDao(handlers ...gdb.ModelHandler) *MediaOutboxEventDao {
	return &MediaOutboxEventDao{
		group:    "default",
		table:    "media_outbox_event",
		columns:  mediaOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MediaOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MediaOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MediaOutboxEventDao) Columns() MediaOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MediaOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MediaOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MediaOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
