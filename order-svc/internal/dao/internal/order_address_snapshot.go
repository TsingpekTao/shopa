// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrderAddressSnapshotDao is the data access object for the table order_address_snapshot.
type OrderAddressSnapshotDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  OrderAddressSnapshotColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// OrderAddressSnapshotColumns defines and stores column names for the table order_address_snapshot.
type OrderAddressSnapshotColumns struct {
	Id                   string //
	OrderNo              string //
	SourceAddressId      string //
	SourceAddressVersion string //
	ReceiverName         string //
	ReceiverPhone        string //
	CountryCode          string //
	ProvinceCode         string //
	ProvinceName         string //
	CityCode             string //
	CityName             string //
	DistrictCode         string //
	DistrictName         string //
	Street               string //
	Detail               string //
	PostalCode           string //
	Latitude             string //
	Longitude            string //
	CreatedAt            string //
	UpdatedAt            string //
}

// orderAddressSnapshotColumns holds the columns for the table order_address_snapshot.
var orderAddressSnapshotColumns = OrderAddressSnapshotColumns{
	Id:                   "id",
	OrderNo:              "order_no",
	SourceAddressId:      "source_address_id",
	SourceAddressVersion: "source_address_version",
	ReceiverName:         "receiver_name",
	ReceiverPhone:        "receiver_phone",
	CountryCode:          "country_code",
	ProvinceCode:         "province_code",
	ProvinceName:         "province_name",
	CityCode:             "city_code",
	CityName:             "city_name",
	DistrictCode:         "district_code",
	DistrictName:         "district_name",
	Street:               "street",
	Detail:               "detail",
	PostalCode:           "postal_code",
	Latitude:             "latitude",
	Longitude:            "longitude",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewOrderAddressSnapshotDao creates and returns a new DAO object for table data access.
func NewOrderAddressSnapshotDao(handlers ...gdb.ModelHandler) *OrderAddressSnapshotDao {
	return &OrderAddressSnapshotDao{
		group:    "default",
		table:    "order_address_snapshot",
		columns:  orderAddressSnapshotColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OrderAddressSnapshotDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OrderAddressSnapshotDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OrderAddressSnapshotDao) Columns() OrderAddressSnapshotColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OrderAddressSnapshotDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OrderAddressSnapshotDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OrderAddressSnapshotDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
