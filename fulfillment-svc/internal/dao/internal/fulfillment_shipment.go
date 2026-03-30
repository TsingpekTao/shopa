// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FulfillmentShipmentDao is the data access object for the table fulfillment_shipment.
type FulfillmentShipmentDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  FulfillmentShipmentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// FulfillmentShipmentColumns defines and stores column names for the table fulfillment_shipment.
type FulfillmentShipmentColumns struct {
	Id                   string //
	ShipmentNo           string //
	OrderNo              string //
	SubOrderNo           string //
	ShopNo               string //
	UserId               string //
	LogisticsCompanyCode string //
	LogisticsCompanyName string //
	LogisticsNo          string //
	ShipmentStatus       string // 1 WAIT_SHIP 2 SHIPPED 3 IN_TRANSIT 4 DELIVERED 5 EXCEPTION 6 CLOSED
	ShippedAt            string //
	DeliveredAt          string //
	ReceiverName         string //
	ReceiverPhone        string //
	ReceiverAddress      string // deprecated flattened address
	ReceiverCountryCode  string //
	ReceiverProvinceCode string //
	ReceiverProvinceName string //
	ReceiverCityCode     string //
	ReceiverCityName     string //
	ReceiverDistrictCode string //
	ReceiverDistrictName string //
	ReceiverStreet       string //
	ReceiverDetail       string //
	ReceiverPostalCode   string //
	Version              string //
	CreatedAt            string //
	UpdatedAt            string //
	DeletedAt            string //
}

// fulfillmentShipmentColumns holds the columns for the table fulfillment_shipment.
var fulfillmentShipmentColumns = FulfillmentShipmentColumns{
	Id:                   "id",
	ShipmentNo:           "shipment_no",
	OrderNo:              "order_no",
	SubOrderNo:           "sub_order_no",
	ShopNo:               "shop_no",
	UserId:               "user_id",
	LogisticsCompanyCode: "logistics_company_code",
	LogisticsCompanyName: "logistics_company_name",
	LogisticsNo:          "logistics_no",
	ShipmentStatus:       "shipment_status",
	ShippedAt:            "shipped_at",
	DeliveredAt:          "delivered_at",
	ReceiverName:         "receiver_name",
	ReceiverPhone:        "receiver_phone",
	ReceiverAddress:      "receiver_address",
	ReceiverCountryCode:  "receiver_country_code",
	ReceiverProvinceCode: "receiver_province_code",
	ReceiverProvinceName: "receiver_province_name",
	ReceiverCityCode:     "receiver_city_code",
	ReceiverCityName:     "receiver_city_name",
	ReceiverDistrictCode: "receiver_district_code",
	ReceiverDistrictName: "receiver_district_name",
	ReceiverStreet:       "receiver_street",
	ReceiverDetail:       "receiver_detail",
	ReceiverPostalCode:   "receiver_postal_code",
	Version:              "version",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
	DeletedAt:            "deleted_at",
}

// NewFulfillmentShipmentDao creates and returns a new DAO object for table data access.
func NewFulfillmentShipmentDao(handlers ...gdb.ModelHandler) *FulfillmentShipmentDao {
	return &FulfillmentShipmentDao{
		group:    "default",
		table:    "fulfillment_shipment",
		columns:  fulfillmentShipmentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FulfillmentShipmentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FulfillmentShipmentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FulfillmentShipmentDao) Columns() FulfillmentShipmentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FulfillmentShipmentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FulfillmentShipmentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FulfillmentShipmentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
