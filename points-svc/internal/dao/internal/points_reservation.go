// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsReservationDao is the data access object for the table points_reservation.
type PointsReservationDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  PointsReservationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// PointsReservationColumns defines and stores column names for the table points_reservation.
type PointsReservationColumns struct {
	ReservationNo         string //
	UserId                string //
	OrderNo               string //
	ShopNo                string //
	ReservationStatusCode string // LOCKED/CONFIRMED/CANCELED/EXPIRED
	RequestedPoints       string //
	LockedPoints          string //
	LockedCashAmountCent  string //
	DeductionDigest       string //
	RuleSnapshotJson      string //
	IdempotencyKey        string //
	ExpireAt              string //
	ConfirmedAt           string //
	CanceledAt            string //
	CancelReasonCode      string //
	CreatedAt             string //
	UpdatedAt             string //
}

// pointsReservationColumns holds the columns for the table points_reservation.
var pointsReservationColumns = PointsReservationColumns{
	ReservationNo:         "reservation_no",
	UserId:                "user_id",
	OrderNo:               "order_no",
	ShopNo:                "shop_no",
	ReservationStatusCode: "reservation_status_code",
	RequestedPoints:       "requested_points",
	LockedPoints:          "locked_points",
	LockedCashAmountCent:  "locked_cash_amount_cent",
	DeductionDigest:       "deduction_digest",
	RuleSnapshotJson:      "rule_snapshot_json",
	IdempotencyKey:        "idempotency_key",
	ExpireAt:              "expire_at",
	ConfirmedAt:           "confirmed_at",
	CanceledAt:            "canceled_at",
	CancelReasonCode:      "cancel_reason_code",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
}

// NewPointsReservationDao creates and returns a new DAO object for table data access.
func NewPointsReservationDao(handlers ...gdb.ModelHandler) *PointsReservationDao {
	return &PointsReservationDao{
		group:    "default",
		table:    "points_reservation",
		columns:  pointsReservationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsReservationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsReservationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsReservationDao) Columns() PointsReservationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsReservationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsReservationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointsReservationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
