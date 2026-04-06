// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsConsumptionDetailDao is the data access object for the table points_consumption_detail.
type PointsConsumptionDetailDao struct {
	table    string                         // table is the underlying table name of the DAO.
	group    string                         // group is the database configuration group name of the current DAO.
	columns  PointsConsumptionDetailColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler             // handlers for customized model modification.
}

// PointsConsumptionDetailColumns defines and stores column names for the table points_consumption_detail.
type PointsConsumptionDetailColumns struct {
	Id             string //
	DetailNo       string //
	UserId         string //
	ReservationNo  string //
	OrderNo        string //
	RefundNo       string //
	SubOrderNo     string //
	BucketNo       string //
	DetailStatusCode string // LOCKED/CONFIRMED/CANCELED
	ConsumedPoints string //
	ReturnedPoints string //
	CashAmountCent string //
	CreatedAt      string //
	UpdatedAt      string //
}

// pointsConsumptionDetailColumns holds the columns for the table points_consumption_detail.
var pointsConsumptionDetailColumns = PointsConsumptionDetailColumns{
	Id:             "id",
	DetailNo:       "detail_no",
	UserId:         "user_id",
	ReservationNo:  "reservation_no",
	OrderNo:        "order_no",
	RefundNo:       "refund_no",
	SubOrderNo:     "sub_order_no",
	BucketNo:       "bucket_no",
	DetailStatusCode: "detail_status_code",
	ConsumedPoints: "consumed_points",
	ReturnedPoints: "returned_points",
	CashAmountCent: "cash_amount_cent",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewPointsConsumptionDetailDao creates and returns a new DAO object for table data access.
func NewPointsConsumptionDetailDao(handlers ...gdb.ModelHandler) *PointsConsumptionDetailDao {
	return &PointsConsumptionDetailDao{
		group:    "default",
		table:    "points_consumption_detail",
		columns:  pointsConsumptionDetailColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsConsumptionDetailDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsConsumptionDetailDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsConsumptionDetailDao) Columns() PointsConsumptionDetailColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsConsumptionDetailDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsConsumptionDetailDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointsConsumptionDetailDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
