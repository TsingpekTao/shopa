// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IamUserAuthDao is the data access object for the table iam_user_auth.
type IamUserAuthDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  IamUserAuthColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// IamUserAuthColumns defines and stores column names for the table iam_user_auth.
type IamUserAuthColumns struct {
	UserId               string // Snowflake user id
	Phone                string // Phone number
	Email                string // Email address
	PasswordHash         string // Password hash
	PasswordSalt         string // Password salt
	PasswordAlgo         string // Password hash algorithm
	PasswordVer          string // Password algorithm version
	AccountStatus        string // 1 active,2 locked,3 disabled,4 review
	FailedLoginCount     string // Fail count for audit
	LockedUntil          string // Account lock expiry time
	TokenVersion         string // Token version for global token invalidation
	LastLoginAt          string // Last login time
	LastLoginIp          string // Last login ip
	LastLoginGeo         string // Last login geo
	LastLoginUa          string // Last login user-agent
	LastLoginFingerprint string // Last login fingerprint
	CreatedAt            string //
	UpdatedAt            string //
}

// iamUserAuthColumns holds the columns for the table iam_user_auth.
var iamUserAuthColumns = IamUserAuthColumns{
	UserId:               "user_id",
	Phone:                "phone",
	Email:                "email",
	PasswordHash:         "password_hash",
	PasswordSalt:         "password_salt",
	PasswordAlgo:         "password_algo",
	PasswordVer:          "password_ver",
	AccountStatus:        "account_status",
	FailedLoginCount:     "failed_login_count",
	LockedUntil:          "locked_until",
	TokenVersion:         "token_version",
	LastLoginAt:          "last_login_at",
	LastLoginIp:          "last_login_ip",
	LastLoginGeo:         "last_login_geo",
	LastLoginUa:          "last_login_ua",
	LastLoginFingerprint: "last_login_fingerprint",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewIamUserAuthDao creates and returns a new DAO object for table data access.
func NewIamUserAuthDao(handlers ...gdb.ModelHandler) *IamUserAuthDao {
	return &IamUserAuthDao{
		group:    "default",
		table:    "iam_user_auth",
		columns:  iamUserAuthColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IamUserAuthDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IamUserAuthDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IamUserAuthDao) Columns() IamUserAuthColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IamUserAuthDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IamUserAuthDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IamUserAuthDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
