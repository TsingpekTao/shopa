// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsAccountDao is the data access object for the table points_account.
type PointsAccountDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  PointsAccountColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// PointsAccountColumns defines and stores column names for the table points_account.
type PointsAccountColumns struct {
	UserId              string // User ID
	AvailableBalance    string // Current available points balance, may be negative when debt exists
	FrozenBalance       string // Currently frozen points balance
	StatusCode          string // ACTIVE/FROZEN/DISABLED
	TotalEarnedPoints   string // Lifetime granted points
	TotalUsedPoints     string // Lifetime confirmed spent points
	TotalExpiredPoints  string // Lifetime expired points
	TotalAdjustedPoints string // Lifetime manual adjustment points
	CreatedAt           string //
	UpdatedAt           string //
}

// pointsAccountColumns holds the columns for the table points_account.
var pointsAccountColumns = PointsAccountColumns{
	UserId:              "user_id",
	AvailableBalance:    "available_balance",
	FrozenBalance:       "frozen_balance",
	StatusCode:          "status_code",
	TotalEarnedPoints:   "total_earned_points",
	TotalUsedPoints:     "total_used_points",
	TotalExpiredPoints:  "total_expired_points",
	TotalAdjustedPoints: "total_adjusted_points",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewPointsAccountDao creates and returns a new DAO object for table data access.
func NewPointsAccountDao(handlers ...gdb.ModelHandler) *PointsAccountDao {
	return &PointsAccountDao{
		group:    "default",
		table:    "points_account",
		columns:  pointsAccountColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsAccountDao) Columns() PointsAccountColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsAccountDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointsAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
