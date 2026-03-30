// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CouponLockDao is the data access object for the table coupon_lock.
type CouponLockDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CouponLockColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CouponLockColumns defines and stores column names for the table coupon_lock.
type CouponLockColumns struct {
	Id            string //
	LockNo        string //
	OrderNo       string //
	UserId        string //
	CouponNosJson string //
	Status        string //
	LockExpireAt  string //
	CreatedAt     string //
	UpdatedAt     string //
}

// couponLockColumns holds the columns for the table coupon_lock.
var couponLockColumns = CouponLockColumns{
	Id:            "id",
	LockNo:        "lock_no",
	OrderNo:       "order_no",
	UserId:        "user_id",
	CouponNosJson: "coupon_nos_json",
	Status:        "status",
	LockExpireAt:  "lock_expire_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewCouponLockDao creates and returns a new DAO object for table data access.
func NewCouponLockDao(handlers ...gdb.ModelHandler) *CouponLockDao {
	return &CouponLockDao{
		group:    "default",
		table:    "coupon_lock",
		columns:  couponLockColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CouponLockDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CouponLockDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CouponLockDao) Columns() CouponLockColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CouponLockDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CouponLockDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CouponLockDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
