// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsRefundActionDao is the data access object for the table points_refund_action.
type PointsRefundActionDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  PointsRefundActionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// PointsRefundActionColumns defines and stores column names for the table points_refund_action.
type PointsRefundActionColumns struct {
	Id                string //
	ActionNo          string //
	ActionType        string // RETURN/REVERSE
	RefundNo          string //
	OrderNo           string //
	SubOrderNo        string //
	ShopNo            string //
	UserId            string //
	IdempotencyKey    string //
	RequestedPoints   string //
	EffectivePoints   string //
	CashAmountCent    string //
	ActionStatus      string // SUCCESS
	ResultPayloadJson string //
	CreatedAt         string //
	UpdatedAt         string //
}

var pointsRefundActionColumns = PointsRefundActionColumns{
	Id:                "id",
	ActionNo:          "action_no",
	ActionType:        "action_type",
	RefundNo:          "refund_no",
	OrderNo:           "order_no",
	SubOrderNo:        "sub_order_no",
	ShopNo:            "shop_no",
	UserId:            "user_id",
	IdempotencyKey:    "idempotency_key",
	RequestedPoints:   "requested_points",
	EffectivePoints:   "effective_points",
	CashAmountCent:    "cash_amount_cent",
	ActionStatus:      "action_status",
	ResultPayloadJson: "result_payload_json",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewPointsRefundActionDao creates and returns a new DAO object for table data access.
func NewPointsRefundActionDao(handlers ...gdb.ModelHandler) *PointsRefundActionDao {
	return &PointsRefundActionDao{
		group:    "default",
		table:    "points_refund_action",
		columns:  pointsRefundActionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsRefundActionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsRefundActionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsRefundActionDao) Columns() PointsRefundActionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsRefundActionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsRefundActionDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *PointsRefundActionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
