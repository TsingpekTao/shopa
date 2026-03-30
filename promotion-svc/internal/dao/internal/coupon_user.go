// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CouponUserDao is the data access object for the table coupon_user.
type CouponUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CouponUserColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CouponUserColumns defines and stores column names for the table coupon_user.
type CouponUserColumns struct {
	Id              string //
	CouponNo        string //
	TemplateNo      string //
	UserId          string //
	Status          string //
	ThresholdAmount string //
	DiscountAmount  string //
	ExpireAt        string //
	CreatedAt       string //
	UpdatedAt       string //
}

// couponUserColumns holds the columns for the table coupon_user.
var couponUserColumns = CouponUserColumns{
	Id:              "id",
	CouponNo:        "coupon_no",
	TemplateNo:      "template_no",
	UserId:          "user_id",
	Status:          "status",
	ThresholdAmount: "threshold_amount",
	DiscountAmount:  "discount_amount",
	ExpireAt:        "expire_at",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewCouponUserDao creates and returns a new DAO object for table data access.
func NewCouponUserDao(handlers ...gdb.ModelHandler) *CouponUserDao {
	return &CouponUserDao{
		group:    "default",
		table:    "coupon_user",
		columns:  couponUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CouponUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CouponUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CouponUserDao) Columns() CouponUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CouponUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CouponUserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CouponUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
