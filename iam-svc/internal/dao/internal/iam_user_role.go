// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserRoleDao is the data access object for the table iam_user_role.
type IamUserRoleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  IamUserRoleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// IamUserRoleColumns defines and stores column names for the table iam_user_role.
type IamUserRoleColumns struct {
	Id        string //
	UserId    string //
	RoleCode  string // 1 customer,2 seller,3 admin,4 cs
	ScopeType string // 1 global,2 shop
	ScopeId   string //
	Status    string // 1 active,2 disabled
	CreatedAt string //
	UpdatedAt string //
}

// iamUserRoleColumns holds the columns for the table iam_user_role.
var iamUserRoleColumns = IamUserRoleColumns{
	Id:        "id",
	UserId:    "user_id",
	RoleCode:  "role_code",
	ScopeType: "scope_type",
	ScopeId:   "scope_id",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewIamUserRoleDao creates and returns a new DAO object for table data access.
func NewIamUserRoleDao(handlers ...gdb.ModelHandler) *IamUserRoleDao {
	return &IamUserRoleDao{
		group:    "default",
		table:    "iam_user_role",
		columns:  iamUserRoleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserRoleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserRoleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserRoleDao) Columns() IamUserRoleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserRoleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserRoleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserRoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
