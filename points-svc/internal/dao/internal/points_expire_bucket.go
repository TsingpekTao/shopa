// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsExpireBucketDao is the data access object for the table points_expire_bucket.
type PointsExpireBucketDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  PointsExpireBucketColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// PointsExpireBucketColumns defines and stores column names for the table points_expire_bucket.
type PointsExpireBucketColumns struct {
	BucketNo        string //
	UserId          string //
	SourceType      string // GRANT/RETURN/REFUND_GRACE/ADJUST
	SourceNo        string //
	SourceVersion   string //
	BucketPeriod    string // YYYYMM or custom window code
	TotalPoints     string //
	RemainingPoints string //
	LockedPoints    string // Currently locked points reserved by active reservations
	UsedPoints      string //
	ExpiredPoints   string //
	ReturnedPoints  string //
	ExpireAt        string //
	BucketStatus    string // ACTIVE/EXPIRED/CLOSED
	CreatedAt       string //
	UpdatedAt       string //
}

// pointsExpireBucketColumns holds the columns for the table points_expire_bucket.
var pointsExpireBucketColumns = PointsExpireBucketColumns{
	BucketNo:        "bucket_no",
	UserId:          "user_id",
	SourceType:      "source_type",
	SourceNo:        "source_no",
	SourceVersion:   "source_version",
	BucketPeriod:    "bucket_period",
	TotalPoints:     "total_points",
	RemainingPoints: "remaining_points",
	LockedPoints:    "locked_points",
	UsedPoints:      "used_points",
	ExpiredPoints:   "expired_points",
	ReturnedPoints:  "returned_points",
	ExpireAt:        "expire_at",
	BucketStatus:    "bucket_status",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewPointsExpireBucketDao creates and returns a new DAO object for table data access.
func NewPointsExpireBucketDao(handlers ...gdb.ModelHandler) *PointsExpireBucketDao {
	return &PointsExpireBucketDao{
		group:    "default",
		table:    "points_expire_bucket",
		columns:  pointsExpireBucketColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsExpireBucketDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsExpireBucketDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsExpireBucketDao) Columns() PointsExpireBucketColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsExpireBucketDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsExpireBucketDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointsExpireBucketDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
