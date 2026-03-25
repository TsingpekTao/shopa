// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserAddressDao is the data access object for the table user_address.
type UserAddressDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserAddressColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserAddressColumns defines and stores column names for the table user_address.
type UserAddressColumns struct {
	AddressId             string // Address ID
	UserId                string // Owner user ID
	Status                string // 1=active,2=deleted,3=replaced
	AddressVersion        string // Optimistic version of address row
	ReplacedFromAddressId string // Address replaced chain source
	Label                 string // Address label
	ReceiverName          string // Receiver name
	ReceiverPhone         string // Receiver phone
	CountryCode           string // ISO country code
	ProvinceCode          string // Province code
	ProvinceName          string // Province name
	CityCode              string // City code
	CityName              string // City name
	DistrictCode          string // District code
	DistrictName          string // District name
	Street                string // Street/town
	Detail                string // Detailed address
	PostalCode            string // Postal code
	IsDefault             string // 1 default, NULL non-default
	DefaultSlot           string // Default uniqueness slot
	Latitude              string // Latitude
	Longitude             string // Longitude
	Ext                   string // Extension payload
	CreatedAt             string //
	UpdatedAt             string //
}

// userAddressColumns holds the columns for the table user_address.
var userAddressColumns = UserAddressColumns{
	AddressId:             "address_id",
	UserId:                "user_id",
	Status:                "status",
	AddressVersion:        "address_version",
	ReplacedFromAddressId: "replaced_from_address_id",
	Label:                 "label",
	ReceiverName:          "receiver_name",
	ReceiverPhone:         "receiver_phone",
	CountryCode:           "country_code",
	ProvinceCode:          "province_code",
	ProvinceName:          "province_name",
	CityCode:              "city_code",
	CityName:              "city_name",
	DistrictCode:          "district_code",
	DistrictName:          "district_name",
	Street:                "street",
	Detail:                "detail",
	PostalCode:            "postal_code",
	IsDefault:             "is_default",
	DefaultSlot:           "default_slot",
	Latitude:              "latitude",
	Longitude:             "longitude",
	Ext:                   "ext",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
}

// NewUserAddressDao creates and returns a new DAO object for table data access.
func NewUserAddressDao(handlers ...gdb.ModelHandler) *UserAddressDao {
	return &UserAddressDao{
		group:    "default",
		table:    "user_address",
		columns:  userAddressColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserAddressDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserAddressDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserAddressDao) Columns() UserAddressColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserAddressDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserAddressDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserAddressDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
