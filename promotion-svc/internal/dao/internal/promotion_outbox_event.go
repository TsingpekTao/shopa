// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionOutboxEventDao is the data access object for the table promotion_outbox_event.
type PromotionOutboxEventDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  PromotionOutboxEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// PromotionOutboxEventColumns defines and stores column names for the table promotion_outbox_event.
type PromotionOutboxEventColumns struct {
	Id          string //
	EventId     string //
	Topic       string //
	EventKey    string //
	PayloadJson string //
	Status      string //
	NextRetryAt string //
	RetryCount  string //
	CreatedAt   string //
	UpdatedAt   string //
}

// promotionOutboxEventColumns holds the columns for the table promotion_outbox_event.
var promotionOutboxEventColumns = PromotionOutboxEventColumns{
	Id:          "id",
	EventId:     "event_id",
	Topic:       "topic",
	EventKey:    "event_key",
	PayloadJson: "payload_json",
	Status:      "status",
	NextRetryAt: "next_retry_at",
	RetryCount:  "retry_count",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewPromotionOutboxEventDao creates and returns a new DAO object for table data access.
func NewPromotionOutboxEventDao(handlers ...gdb.ModelHandler) *PromotionOutboxEventDao {
	return &PromotionOutboxEventDao{
		group:    "default",
		table:    "promotion_outbox_event",
		columns:  promotionOutboxEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PromotionOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PromotionOutboxEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PromotionOutboxEventDao) Columns() PromotionOutboxEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PromotionOutboxEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PromotionOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PromotionOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
