// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsLedgerDao is the data access object for the table points_ledger.
type PointsLedgerDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  PointsLedgerColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// PointsLedgerColumns defines and stores column names for the table points_ledger.
type PointsLedgerColumns struct {
	Id              string //
	LedgerNo        string // Business ledger identifier
	UserId          string //
	EntryTypeCode   string // INIT/LOCK/CONFIRM/CANCEL/GRANT/RETURN/REVERSE/EXPIRE/ADJUST/FREEZE/UNFREEZE
	BizType         string // REGISTER_INIT/ORDER_PAY/REFUND/etc
	BizNo           string // Business identifier for idempotency
	ReservationNo   string // Reservation identifier if relevant
	RelatedBucketNo string // Bucket identifier if relevant
	PointsDelta     string // Points delta for this ledger entry
	AvailableAfter  string // Available balance snapshot after apply
	FrozenAfter     string // Frozen balance snapshot after apply
	DebtAfter       string // Debt snapshot after apply
	CashAmountCent  string // Related cash amount in cents
	Remark          string // Operator remark or domain explanation
	ExtraJson       string // Extended metadata snapshot
	CreatedAt       string //
}

// pointsLedgerColumns holds the columns for the table points_ledger.
var pointsLedgerColumns = PointsLedgerColumns{
	Id:              "id",
	LedgerNo:        "ledger_no",
	UserId:          "user_id",
	EntryTypeCode:   "entry_type_code",
	BizType:         "biz_type",
	BizNo:           "biz_no",
	ReservationNo:   "reservation_no",
	RelatedBucketNo: "related_bucket_no",
	PointsDelta:     "points_delta",
	AvailableAfter:  "available_after",
	FrozenAfter:     "frozen_after",
	DebtAfter:       "debt_after",
	CashAmountCent:  "cash_amount_cent",
	Remark:          "remark",
	ExtraJson:       "extra_json",
	CreatedAt:       "created_at",
}

// NewPointsLedgerDao creates and returns a new DAO object for table data access.
func NewPointsLedgerDao(handlers ...gdb.ModelHandler) *PointsLedgerDao {
	return &PointsLedgerDao{
		group:    "default",
		table:    "points_ledger",
		columns:  pointsLedgerColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsLedgerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsLedgerDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsLedgerDao) Columns() PointsLedgerColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsLedgerDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsLedgerDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointsLedgerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
