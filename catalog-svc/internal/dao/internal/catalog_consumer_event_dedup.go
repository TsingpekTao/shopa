// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogConsumerEventDedupDao is the data access object for the table catalog_consumer_event_dedup.
type CatalogConsumerEventDedupDao struct {
	table    string                           // table is the underlying table name of the DAO.
	group    string                           // group is the database configuration group name of the current DAO.
	columns  CatalogConsumerEventDedupColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler               // handlers for customized model modification.
}

// CatalogConsumerEventDedupColumns defines and stores column names for the table catalog_consumer_event_dedup.
type CatalogConsumerEventDedupColumns struct {
	Id           string //
	ConsumerName string //
	EventId      string //
	PayloadHash  string //
	FirstSeenAt  string //
}

// catalogConsumerEventDedupColumns holds the columns for the table catalog_consumer_event_dedup.
var catalogConsumerEventDedupColumns = CatalogConsumerEventDedupColumns{
	Id:           "id",
	ConsumerName: "consumer_name",
	EventId:      "event_id",
	PayloadHash:  "payload_hash",
	FirstSeenAt:  "first_seen_at",
}

// NewCatalogConsumerEventDedupDao creates and returns a new DAO object for table data access.
func NewCatalogConsumerEventDedupDao(handlers ...gdb.ModelHandler) *CatalogConsumerEventDedupDao {
	return &CatalogConsumerEventDedupDao{
		group:    "default",
		table:    "catalog_consumer_event_dedup",
		columns:  catalogConsumerEventDedupColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogConsumerEventDedupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogConsumerEventDedupDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogConsumerEventDedupDao) Columns() CatalogConsumerEventDedupColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogConsumerEventDedupDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogConsumerEventDedupDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogConsumerEventDedupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
