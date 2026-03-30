// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentReconciliationRecordDao is the data access object for the table payment_reconciliation_record.
type PaymentReconciliationRecordDao struct {
	table    string                             // table is the underlying table name of the DAO.
	group    string                             // group is the database configuration group name of the current DAO.
	columns  PaymentReconciliationRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                 // handlers for customized model modification.
}

// PaymentReconciliationRecordColumns defines and stores column names for the table payment_reconciliation_record.
type PaymentReconciliationRecordColumns struct {
	Id            string //
	DiffNo        string //
	ReconTaskNo   string //
	DiffType      string //
	PaymentNo     string //
	OrderNo       string //
	LocalAmount   string //
	GatewayAmount string //
	Status        string //
	DetailJson    string //
	ResolvedAt    string //
	CreatedAt     string //
	UpdatedAt     string //
}

// paymentReconciliationRecordColumns holds the columns for the table payment_reconciliation_record.
var paymentReconciliationRecordColumns = PaymentReconciliationRecordColumns{
	Id:            "id",
	DiffNo:        "diff_no",
	ReconTaskNo:   "recon_task_no",
	DiffType:      "diff_type",
	PaymentNo:     "payment_no",
	OrderNo:       "order_no",
	LocalAmount:   "local_amount",
	GatewayAmount: "gateway_amount",
	Status:        "status",
	DetailJson:    "detail_json",
	ResolvedAt:    "resolved_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewPaymentReconciliationRecordDao creates and returns a new DAO object for table data access.
func NewPaymentReconciliationRecordDao(handlers ...gdb.ModelHandler) *PaymentReconciliationRecordDao {
	return &PaymentReconciliationRecordDao{
		group:    "default",
		table:    "payment_reconciliation_record",
		columns:  paymentReconciliationRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PaymentReconciliationRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PaymentReconciliationRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PaymentReconciliationRecordDao) Columns() PaymentReconciliationRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PaymentReconciliationRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PaymentReconciliationRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PaymentReconciliationRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
