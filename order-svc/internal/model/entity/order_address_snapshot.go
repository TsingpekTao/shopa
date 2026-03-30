// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderAddressSnapshot is the golang structure for table order_address_snapshot.
type OrderAddressSnapshot struct {
	Id                   uint64      `json:"id"                   orm:"id"                     description:""` //
	OrderNo              string      `json:"orderNo"              orm:"order_no"               description:""` //
	SourceAddressId      uint64      `json:"sourceAddressId"      orm:"source_address_id"      description:""` //
	SourceAddressVersion uint64      `json:"sourceAddressVersion" orm:"source_address_version" description:""` //
	ReceiverName         string      `json:"receiverName"         orm:"receiver_name"          description:""` //
	ReceiverPhone        string      `json:"receiverPhone"        orm:"receiver_phone"         description:""` //
	CountryCode          string      `json:"countryCode"          orm:"country_code"           description:""` //
	ProvinceCode         string      `json:"provinceCode"         orm:"province_code"          description:""` //
	ProvinceName         string      `json:"provinceName"         orm:"province_name"          description:""` //
	CityCode             string      `json:"cityCode"             orm:"city_code"              description:""` //
	CityName             string      `json:"cityName"             orm:"city_name"              description:""` //
	DistrictCode         string      `json:"districtCode"         orm:"district_code"          description:""` //
	DistrictName         string      `json:"districtName"         orm:"district_name"          description:""` //
	Street               string      `json:"street"               orm:"street"                 description:""` //
	Detail               string      `json:"detail"               orm:"detail"                 description:""` //
	PostalCode           string      `json:"postalCode"           orm:"postal_code"            description:""` //
	Latitude             float64     `json:"latitude"             orm:"latitude"               description:""` //
	Longitude            float64     `json:"longitude"            orm:"longitude"              description:""` //
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             description:""` //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:""` //
}
