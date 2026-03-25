// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserOauthDao is the data access object for the table iam_user_oauth.
type IamUserOauthDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  IamUserOauthColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// IamUserOauthColumns defines and stores column names for the table iam_user_oauth.
type IamUserOauthColumns struct {
	Id          string //
	UserId      string //
	Provider    string // 1 wechat,2 alipay
	ProviderUid string //
	UnionId     string //
	Status      string // 1 active,2 unbound
	CreatedAt   string //
	UpdatedAt   string //
}

// iamUserOauthColumns holds the columns for the table iam_user_oauth.
var iamUserOauthColumns = IamUserOauthColumns{
	Id:          "id",
	UserId:      "user_id",
	Provider:    "provider",
	ProviderUid: "provider_uid",
	UnionId:     "union_id",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewIamUserOauthDao creates and returns a new DAO object for table data access.
func NewIamUserOauthDao(handlers ...gdb.ModelHandler) *IamUserOauthDao {
	return &IamUserOauthDao{
		group:    "default",
		table:    "iam_user_oauth",
		columns:  iamUserOauthColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserOauthDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserOauthDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserOauthDao) Columns() IamUserOauthColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserOauthDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserOauthDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserOauthDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
