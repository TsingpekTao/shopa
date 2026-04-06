// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderAddressSnapshot is the golang structure for table order_address_snapshot.
type OrderAddressSnapshot struct {
	Id                   uint64      `json:"id"                   orm:"id"                     ` //
	OrderNo              string      `json:"orderNo"              orm:"order_no"               ` //
	SourceAddressId      uint64      `json:"sourceAddressId"      orm:"source_address_id"      ` //
	SourceAddressVersion uint64      `json:"sourceAddressVersion" orm:"source_address_version" ` //
	ReceiverName         string      `json:"receiverName"         orm:"receiver_name"          ` //
	ReceiverPhone        string      `json:"receiverPhone"        orm:"receiver_phone"         ` //
	CountryCode          string      `json:"countryCode"          orm:"country_code"           ` //
	ProvinceCode         string      `json:"provinceCode"         orm:"province_code"          ` //
	ProvinceName         string      `json:"provinceName"         orm:"province_name"          ` //
	CityCode             string      `json:"cityCode"             orm:"city_code"              ` //
	CityName             string      `json:"cityName"             orm:"city_name"              ` //
	DistrictCode         string      `json:"districtCode"         orm:"district_code"          ` //
	DistrictName         string      `json:"districtName"         orm:"district_name"          ` //
	Street               string      `json:"street"               orm:"street"                 ` //
	Detail               string      `json:"detail"               orm:"detail"                 ` //
	PostalCode           string      `json:"postalCode"           orm:"postal_code"            ` //
	Latitude             float64     `json:"latitude"             orm:"latitude"               ` //
	Longitude            float64     `json:"longitude"            orm:"longitude"              ` //
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             ` //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             ` //
}
