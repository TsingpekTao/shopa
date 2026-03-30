// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FulfillmentTrackingNodeDao is the data access object for the table fulfillment_tracking_node.
type FulfillmentTrackingNodeDao struct {
	table    string                         // table is the underlying table name of the DAO.
	group    string                         // group is the database configuration group name of the current DAO.
	columns  FulfillmentTrackingNodeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler             // handlers for customized model modification.
}

// FulfillmentTrackingNodeColumns defines and stores column names for the table fulfillment_tracking_node.
type FulfillmentTrackingNodeColumns struct {
	Id         string //
	NodeNo     string //
	ShipmentNo string //
	StatusCode string //
	Content    string //
	Location   string //
	EventTime  string //
	CreatedAt  string //
}

// fulfillmentTrackingNodeColumns holds the columns for the table fulfillment_tracking_node.
var fulfillmentTrackingNodeColumns = FulfillmentTrackingNodeColumns{
	Id:         "id",
	NodeNo:     "node_no",
	ShipmentNo: "shipment_no",
	StatusCode: "status_code",
	Content:    "content",
	Location:   "location",
	EventTime:  "event_time",
	CreatedAt:  "created_at",
}

// NewFulfillmentTrackingNodeDao creates and returns a new DAO object for table data access.
func NewFulfillmentTrackingNodeDao(handlers ...gdb.ModelHandler) *FulfillmentTrackingNodeDao {
	return &FulfillmentTrackingNodeDao{
		group:    "default",
		table:    "fulfillment_tracking_node",
		columns:  fulfillmentTrackingNodeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FulfillmentTrackingNodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FulfillmentTrackingNodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FulfillmentTrackingNodeDao) Columns() FulfillmentTrackingNodeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FulfillmentTrackingNodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FulfillmentTrackingNodeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FulfillmentTrackingNodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
