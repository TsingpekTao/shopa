// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsGrantDetailDao is the data access object for the table points_grant_detail.
type PointsGrantDetailDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  PointsGrantDetailColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// PointsGrantDetailColumns defines and stores column names for the table points_grant_detail.
type PointsGrantDetailColumns struct {
	Id              string //
	GrantDetailNo   string //
	UserId          string //
	OrderNo         string //
	SubOrderNo      string //
	ShopNo          string //
	GrantedPoints   string //
	ReversedPoints  string //
	RuleCode        string //
	RuleSnapshotJson string //
	GrantStatusCode string // GRANTED/PARTIAL_REVERSED/REVERSED
	CreatedAt       string //
	UpdatedAt       string //
}

// pointsGrantDetailColumns holds the columns for the table points_grant_detail.
var pointsGrantDetailColumns = PointsGrantDetailColumns{
	Id:              "id",
	GrantDetailNo:   "grant_detail_no",
	UserId:          "user_id",
	OrderNo:         "order_no",
	SubOrderNo:      "sub_order_no",
	ShopNo:          "shop_no",
	GrantedPoints:   "granted_points",
	ReversedPoints:  "reversed_points",
	RuleCode:        "rule_code",
	RuleSnapshotJson: "rule_snapshot_json",
	GrantStatusCode: "grant_status_code",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewPointsGrantDetailDao creates and returns a new DAO object for table data access.
func NewPointsGrantDetailDao(handlers ...gdb.ModelHandler) *PointsGrantDetailDao {
	return &PointsGrantDetailDao{
		group:    "default",
		table:    "points_grant_detail",
		columns:  pointsGrantDetailColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsGrantDetailDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsGrantDetailDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsGrantDetailDao) Columns() PointsGrantDetailColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsGrantDetailDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsGrantDetailDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointsGrantDetailDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
