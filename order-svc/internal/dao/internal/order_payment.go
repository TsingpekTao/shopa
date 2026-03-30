// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderPaymentDao is the data access object for the table order_payment.
type OrderPaymentDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  OrderPaymentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// OrderPaymentColumns defines and stores column names for the table order_payment.
type OrderPaymentColumns struct {
	Id             string //
	OrderNo        string //
	PayNo          string //
	PaymentEventId string //
	PayChannel     string //
	PayStatusCode  string //
	ChannelTradeNo string //
	PaidAmount     string //
	PaidAt         string //
	RawPayload     string //
	CreatedAt      string //
	UpdatedAt      string //
}

// orderPaymentColumns holds the columns for the table order_payment.
var orderPaymentColumns = OrderPaymentColumns{
	Id:             "id",
	OrderNo:        "order_no",
	PayNo:          "pay_no",
	PaymentEventId: "payment_event_id",
	PayChannel:     "pay_channel",
	PayStatusCode:  "pay_status_code",
	ChannelTradeNo: "channel_trade_no",
	PaidAmount:     "paid_amount",
	PaidAt:         "paid_at",
	RawPayload:     "raw_payload",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewOrderPaymentDao creates and returns a new DAO object for table data access.
func NewOrderPaymentDao(handlers ...gdb.ModelHandler) *OrderPaymentDao {
	return &OrderPaymentDao{
		group:    "default",
		table:    "order_payment",
		columns:  orderPaymentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrderPaymentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrderPaymentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrderPaymentDao) Columns() OrderPaymentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrderPaymentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrderPaymentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OrderPaymentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
