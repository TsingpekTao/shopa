// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentCallbackLogDao is the data access object for the table payment_callback_log.
type PaymentCallbackLogDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  PaymentCallbackLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// PaymentCallbackLogColumns defines and stores column names for the table payment_callback_log.
type PaymentCallbackLogColumns struct {
	Id                string //
	CallbackEventId   string //
	PaymentNo         string //
	OrderNo           string //
	PayChannel        string //
	GatewayStatusCode string //
	PaidAmount        string //
	RawPayload        string //
	Signature         string //
	CreatedAt         string //
}

// paymentCallbackLogColumns holds the columns for the table payment_callback_log.
var paymentCallbackLogColumns = PaymentCallbackLogColumns{
	Id:                "id",
	CallbackEventId:   "callback_event_id",
	PaymentNo:         "payment_no",
	OrderNo:           "order_no",
	PayChannel:        "pay_channel",
	GatewayStatusCode: "gateway_status_code",
	PaidAmount:        "paid_amount",
	RawPayload:        "raw_payload",
	Signature:         "signature",
	CreatedAt:         "created_at",
}

// NewPaymentCallbackLogDao creates and returns a new DAO object for table data access.
func NewPaymentCallbackLogDao(handlers ...gdb.ModelHandler) *PaymentCallbackLogDao {
	return &PaymentCallbackLogDao{
		group:    "default",
		table:    "payment_callback_log",
		columns:  paymentCallbackLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PaymentCallbackLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PaymentCallbackLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PaymentCallbackLogDao) Columns() PaymentCallbackLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PaymentCallbackLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PaymentCallbackLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PaymentCallbackLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
