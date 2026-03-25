// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SellerApplicationOpLogDao is the data access object for the table seller_application_op_log.
type SellerApplicationOpLogDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  SellerApplicationOpLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// SellerApplicationOpLogColumns defines and stores column names for the table seller_application_op_log.
type SellerApplicationOpLogColumns struct {
	Id             string //
	ApplicationNo  string //
	OwnerUserId    string //
	OperatorUserId string // Admin or seller operator
	ActionCode     string // CREATE_DRAFT/SUBMIT/APPROVE/REJECT/RESUBMIT
	FromStatus     string //
	ToStatus       string //
	Comment        string //
	ExtraJson      string //
	CreatedAt      string //
}

// sellerApplicationOpLogColumns holds the columns for the table seller_application_op_log.
var sellerApplicationOpLogColumns = SellerApplicationOpLogColumns{
	Id:             "id",
	ApplicationNo:  "application_no",
	OwnerUserId:    "owner_user_id",
	OperatorUserId: "operator_user_id",
	ActionCode:     "action_code",
	FromStatus:     "from_status",
	ToStatus:       "to_status",
	Comment:        "comment",
	ExtraJson:      "extra_json",
	CreatedAt:      "created_at",
}

// NewSellerApplicationOpLogDao creates and returns a new DAO object for table data access.
func NewSellerApplicationOpLogDao(handlers ...gdb.ModelHandler) *SellerApplicationOpLogDao {
	return &SellerApplicationOpLogDao{
		group:    "default",
		table:    "seller_application_op_log",
		columns:  sellerApplicationOpLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SellerApplicationOpLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SellerApplicationOpLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SellerApplicationOpLogDao) Columns() SellerApplicationOpLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SellerApplicationOpLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SellerApplicationOpLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SellerApplicationOpLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
