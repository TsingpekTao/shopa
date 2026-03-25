// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamRefreshSessionDao is the data access object for the table iam_refresh_session.
type IamRefreshSessionDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  IamRefreshSessionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// IamRefreshSessionColumns defines and stores column names for the table iam_refresh_session.
type IamRefreshSessionColumns struct {
	Sid              string // Session id
	UserId           string //
	RefreshTokenHash string //
	UaHash           string //
	Ip               string //
	Geo              string //
	Fingerprint      string //
	CreatedAt        string //
	ExpiresAt        string //
	RevokedAt        string //
	ReplacedBySid    string //
}

// iamRefreshSessionColumns holds the columns for the table iam_refresh_session.
var iamRefreshSessionColumns = IamRefreshSessionColumns{
	Sid:              "sid",
	UserId:           "user_id",
	RefreshTokenHash: "refresh_token_hash",
	UaHash:           "ua_hash",
	Ip:               "ip",
	Geo:              "geo",
	Fingerprint:      "fingerprint",
	CreatedAt:        "created_at",
	ExpiresAt:        "expires_at",
	RevokedAt:        "revoked_at",
	ReplacedBySid:    "replaced_by_sid",
}

// NewIamRefreshSessionDao creates and returns a new DAO object for table data access.
func NewIamRefreshSessionDao(handlers ...gdb.ModelHandler) *IamRefreshSessionDao {
	return &IamRefreshSessionDao{
		group:    "default",
		table:    "iam_refresh_session",
		columns:  iamRefreshSessionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamRefreshSessionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamRefreshSessionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamRefreshSessionDao) Columns() IamRefreshSessionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamRefreshSessionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamRefreshSessionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamRefreshSessionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
