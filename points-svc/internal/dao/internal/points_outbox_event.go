// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsOutboxEventDao is the data access object for the table points_outbox_event.
type PointsOutboxEventDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  PointsOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// PointsOutboxEventColumns defines and stores column names for the table points_outbox_event.
type PointsOutboxEventColumns struct {
	Id                string //
	EventNo           string //
	AggregateType     string //
	AggregateNo       string //
	EventType         string //
	PayloadJson       string //
	PublishStatusCode string // PENDING/SENT/FAILED/DEAD
	RetryCount        string //
	NextRetryAt       string //
	LastError         string //
	CreatedAt         string //
	UpdatedAt         string //
}

// pointsOutboxEventColumns holds the columns for the table points_outbox_event.
var pointsOutboxEventColumns = PointsOutboxEventColumns{
	Id:                "id",
	EventNo:           "event_no",
	AggregateType:     "aggregate_type",
	AggregateNo:       "aggregate_no",
	EventType:         "event_type",
	PayloadJson:       "payload_json",
	PublishStatusCode: "publish_status_code",
	RetryCount:        "retry_count",
	NextRetryAt:       "next_retry_at",
	LastError:         "last_error",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewPointsOutboxEventDao creates and returns a new DAO object for table data access.
func NewPointsOutboxEventDao(handlers ...gdb.ModelHandler) *PointsOutboxEventDao {
	return &PointsOutboxEventDao{
		group:    "default",
		table:    "points_outbox_event",
		columns:  pointsOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsOutboxEventDao) Columns() PointsOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointsOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
