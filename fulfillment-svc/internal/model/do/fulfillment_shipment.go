// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentShipment is the golang structure of table fulfillment_shipment for DAO operations like Where/Data.
type FulfillmentShipment struct {
	g.Meta               `orm:"table:fulfillment_shipment, do:true"`
	Id                   any         //
	ShipmentNo           any         //
	OrderNo              any         //
	SubOrderNo           any         //
	ShopNo               any         //
	UserId               any         //
	LogisticsCompanyCode any         //
	LogisticsCompanyName any         //
	LogisticsNo          any         //
	ShipmentStatus       any         // 1 WAIT_SHIP 2 SHIPPED 3 IN_TRANSIT 4 DELIVERED 5 EXCEPTION 6 CLOSED
	ShippedAt            *gtime.Time //
	DeliveredAt          *gtime.Time //
	ReceiverName         any         //
	ReceiverPhone        any         //
	ReceiverAddress      any         // deprecated flattened address
	ReceiverCountryCode  any         //
	ReceiverProvinceCode any         //
	ReceiverProvinceName any         //
	ReceiverCityCode     any         //
	ReceiverCityName     any         //
	ReceiverDistrictCode any         //
	ReceiverDistrictName any         //
	ReceiverStreet       any         //
	ReceiverDetail       any         //
	ReceiverPostalCode   any         //
	Version              any         //
	CreatedAt            *gtime.Time //
	UpdatedAt            *gtime.Time //
	DeletedAt            *gtime.Time //
}
