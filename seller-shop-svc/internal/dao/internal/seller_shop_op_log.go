// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SellerShopOpLogDao is the data access object for the table seller_shop_op_log.
type SellerShopOpLogDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  SellerShopOpLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// SellerShopOpLogColumns defines and stores column names for the table seller_shop_op_log.
type SellerShopOpLogColumns struct {
	Id             string //
	ShopNo         string //
	OwnerUserId    string //
	OperatorUserId string //
	ActionCode     string // PROVISION_START/ACTIVATE/FREEZE/CLOSE/REOPEN
	FromStatus     string //
	ToStatus       string //
	ReasonCode     string //
	Reason         string //
	RequestId      string //
	EventId        string //
	CreatedAt      string //
}

// sellerShopOpLogColumns holds the columns for the table seller_shop_op_log.
var sellerShopOpLogColumns = SellerShopOpLogColumns{
	Id:             "id",
	ShopNo:         "shop_no",
	OwnerUserId:    "owner_user_id",
	OperatorUserId: "operator_user_id",
	ActionCode:     "action_code",
	FromStatus:     "from_status",
	ToStatus:       "to_status",
	ReasonCode:     "reason_code",
	Reason:         "reason",
	RequestId:      "request_id",
	EventId:        "event_id",
	CreatedAt:      "created_at",
}

// NewSellerShopOpLogDao creates and returns a new DAO object for table data access.
func NewSellerShopOpLogDao(handlers ...gdb.ModelHandler) *SellerShopOpLogDao {
	return &SellerShopOpLogDao{
		group:    "default",
		table:    "seller_shop_op_log",
		columns:  sellerShopOpLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SellerShopOpLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SellerShopOpLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SellerShopOpLogDao) Columns() SellerShopOpLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SellerShopOpLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SellerShopOpLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SellerShopOpLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
