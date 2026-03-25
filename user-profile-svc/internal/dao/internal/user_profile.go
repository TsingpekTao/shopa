// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserProfileDao is the data access object for the table user_profile.
type UserProfileDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserProfileColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserProfileColumns defines and stores column names for the table user_profile.
type UserProfileColumns struct {
	UserId             string // User ID from iam-svc
	DisplayName        string // Display name
	AvatarAssetId      string // Avatar asset id from media-svc
	AvatarUrl          string // Avatar URL cache
	Gender             string // 0=unspecified,1=male,2=female,3=other
	Birthday           string // Birthday date
	Locale             string // Locale
	Timezone           string // Timezone
	MarketingOptIn     string // Marketing opt in
	DefaultAddressId   string // Deprecated compatibility field
	ProfileVersion     string // Optimistic version of profile aggregate
	AddressBookVersion string // CAS version of address-book aggregate
	DisplayNameSource  string // 1 SYSTEM_INIT,2 USER_SET
	LastInitEventAt    string // Last register-init event time
	LastInitEventId    string // Last register-init event id
	Ext                string // Extension payload
	CreatedAt          string //
	UpdatedAt          string //
}

// userProfileColumns holds the columns for the table user_profile.
var userProfileColumns = UserProfileColumns{
	UserId:             "user_id",
	DisplayName:        "display_name",
	AvatarAssetId:      "avatar_asset_id",
	AvatarUrl:          "avatar_url",
	Gender:             "gender",
	Birthday:           "birthday",
	Locale:             "locale",
	Timezone:           "timezone",
	MarketingOptIn:     "marketing_opt_in",
	DefaultAddressId:   "default_address_id",
	ProfileVersion:     "profile_version",
	AddressBookVersion: "address_book_version",
	DisplayNameSource:  "display_name_source",
	LastInitEventAt:    "last_init_event_at",
	LastInitEventId:    "last_init_event_id",
	Ext:                "ext",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
}

// NewUserProfileDao creates and returns a new DAO object for table data access.
func NewUserProfileDao(handlers ...gdb.ModelHandler) *UserProfileDao {
	return &UserProfileDao{
		group:    "default",
		table:    "user_profile",
		columns:  userProfileColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserProfileDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserProfileDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserProfileDao) Columns() UserProfileColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserProfileDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserProfileDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserProfileDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
