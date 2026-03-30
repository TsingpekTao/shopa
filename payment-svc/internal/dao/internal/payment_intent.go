// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentIntentDao is the data access object for the table payment_intent.
type PaymentIntentDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  PaymentIntentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// PaymentIntentColumns defines and stores column names for the table payment_intent.
type PaymentIntentColumns struct {
	Id              string //
	PaymentNo       string //
	OrderNo         string //
	UserId          string //
	PayChannel      string //
	Status          string //
	PayableAmount   string //
	RefundedAmount  string //
	CurrencyCode    string //
	OrderExpireAt   string //
	GatewayExpireAt string //
	ExternalTradeNo string //
	PaidAt          string //
	ClosedAt        string //
	Version         string //
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string //
}

// paymentIntentColumns holds the columns for the table payment_intent.
var paymentIntentColumns = PaymentIntentColumns{
	Id:              "id",
	PaymentNo:       "payment_no",
	OrderNo:         "order_no",
	UserId:          "user_id",
	PayChannel:      "pay_channel",
	Status:          "status",
	PayableAmount:   "payable_amount",
	RefundedAmount:  "refunded_amount",
	CurrencyCode:    "currency_code",
	OrderExpireAt:   "order_expire_at",
	GatewayExpireAt: "gateway_expire_at",
	ExternalTradeNo: "external_trade_no",
	PaidAt:          "paid_at",
	ClosedAt:        "closed_at",
	Version:         "version",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewPaymentIntentDao creates and returns a new DAO object for table data access.
func NewPaymentIntentDao(handlers ...gdb.ModelHandler) *PaymentIntentDao {
	return &PaymentIntentDao{
		group:    "default",
		table:    "payment_intent",
		columns:  paymentIntentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PaymentIntentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PaymentIntentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PaymentIntentDao) Columns() PaymentIntentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PaymentIntentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PaymentIntentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PaymentIntentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
