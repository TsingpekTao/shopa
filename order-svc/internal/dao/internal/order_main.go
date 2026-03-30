// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderMainDao is the data access object for the table order_main.
type OrderMainDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  OrderMainColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// OrderMainColumns defines and stores column names for the table order_main.
type OrderMainColumns struct {
	Id               string //
	OrderNo          string //
	UserId           string //
	OrderStatus      string //
	PaymentStatus    string //
	ReservationNo    string //
	GoodsAmount      string //
	FreightAmount    string //
	DiscountAmount   string //
	PayableAmount    string //
	PaidAmount       string //
	BuyerRemark      string //
	CancelReasonCode string //
	PayDeadlineAt    string //
	PaidAt           string //
	ClosedAt         string //
	Version          string //
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// orderMainColumns holds the columns for the table order_main.
var orderMainColumns = OrderMainColumns{
	Id:               "id",
	OrderNo:          "order_no",
	UserId:           "user_id",
	OrderStatus:      "order_status",
	PaymentStatus:    "payment_status",
	ReservationNo:    "reservation_no",
	GoodsAmount:      "goods_amount",
	FreightAmount:    "freight_amount",
	DiscountAmount:   "discount_amount",
	PayableAmount:    "payable_amount",
	PaidAmount:       "paid_amount",
	BuyerRemark:      "buyer_remark",
	CancelReasonCode: "cancel_reason_code",
	PayDeadlineAt:    "pay_deadline_at",
	PaidAt:           "paid_at",
	ClosedAt:         "closed_at",
	Version:          "version",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewOrderMainDao creates and returns a new DAO object for table data access.
func NewOrderMainDao(handlers ...gdb.ModelHandler) *OrderMainDao {
	return &OrderMainDao{
		group:    "default",
		table:    "order_main",
		columns:  orderMainColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrderMainDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrderMainDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrderMainDao) Columns() OrderMainColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrderMainDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrderMainDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OrderMainDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
