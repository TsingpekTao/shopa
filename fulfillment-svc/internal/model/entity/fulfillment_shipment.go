// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentShipment is the golang structure for table fulfillment_shipment.
type FulfillmentShipment struct {
	Id                   uint64      `json:"id"                   orm:"id"                     description:""`                                                                    //
	ShipmentNo           string      `json:"shipmentNo"           orm:"shipment_no"            description:""`                                                                    //
	OrderNo              string      `json:"orderNo"              orm:"order_no"               description:""`                                                                    //
	SubOrderNo           string      `json:"subOrderNo"           orm:"sub_order_no"           description:""`                                                                    //
	ShopNo               string      `json:"shopNo"               orm:"shop_no"                description:""`                                                                    //
	UserId               uint64      `json:"userId"               orm:"user_id"                description:""`                                                                    //
	LogisticsCompanyCode string      `json:"logisticsCompanyCode" orm:"logistics_company_code" description:""`                                                                    //
	LogisticsCompanyName string      `json:"logisticsCompanyName" orm:"logistics_company_name" description:""`                                                                    //
	LogisticsNo          string      `json:"logisticsNo"          orm:"logistics_no"           description:""`                                                                    //
	ShipmentStatus       uint        `json:"shipmentStatus"       orm:"shipment_status"        description:"1 WAIT_SHIP 2 SHIPPED 3 IN_TRANSIT 4 DELIVERED 5 EXCEPTION 6 CLOSED"` // 1 WAIT_SHIP 2 SHIPPED 3 IN_TRANSIT 4 DELIVERED 5 EXCEPTION 6 CLOSED
	ShippedAt            *gtime.Time `json:"shippedAt"            orm:"shipped_at"             description:""`                                                                    //
	DeliveredAt          *gtime.Time `json:"deliveredAt"          orm:"delivered_at"           description:""`                                                                    //
	ReceiverName         string      `json:"receiverName"         orm:"receiver_name"          description:""`                                                                    //
	ReceiverPhone        string      `json:"receiverPhone"        orm:"receiver_phone"         description:""`                                                                    //
	ReceiverAddress      string      `json:"receiverAddress"      orm:"receiver_address"       description:"deprecated flattened address"`                                        // deprecated flattened address
	ReceiverCountryCode  string      `json:"receiverCountryCode"  orm:"receiver_country_code"  description:""`                                                                    //
	ReceiverProvinceCode string      `json:"receiverProvinceCode" orm:"receiver_province_code" description:""`                                                                    //
	ReceiverProvinceName string      `json:"receiverProvinceName" orm:"receiver_province_name" description:""`                                                                    //
	ReceiverCityCode     string      `json:"receiverCityCode"     orm:"receiver_city_code"     description:""`                                                                    //
	ReceiverCityName     string      `json:"receiverCityName"     orm:"receiver_city_name"     description:""`                                                                    //
	ReceiverDistrictCode string      `json:"receiverDistrictCode" orm:"receiver_district_code" description:""`                                                                    //
	ReceiverDistrictName string      `json:"receiverDistrictName" orm:"receiver_district_name" description:""`                                                                    //
	ReceiverStreet       string      `json:"receiverStreet"       orm:"receiver_street"        description:""`                                                                    //
	ReceiverDetail       string      `json:"receiverDetail"       orm:"receiver_detail"        description:""`                                                                    //
	ReceiverPostalCode   string      `json:"receiverPostalCode"   orm:"receiver_postal_code"   description:""`                                                                    //
	Version              uint64      `json:"version"              orm:"version"                description:""`                                                                    //
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             description:""`                                                                    //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:""`                                                                    //
	DeletedAt            *gtime.Time `json:"deletedAt"            orm:"deleted_at"             description:""`                                                                    //
}
