// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAddress is the golang structure for table user_address.
type UserAddress struct {
	AddressId             uint64      `json:"addressId"             orm:"address_id"               description:"Address ID"`                        // Address ID
	UserId                uint64      `json:"userId"                orm:"user_id"                  description:"Owner user ID"`                     // Owner user ID
	Status                uint        `json:"status"                orm:"status"                   description:"1=active,2=deleted,3=replaced"`     // 1=active,2=deleted,3=replaced
	AddressVersion        uint64      `json:"addressVersion"        orm:"address_version"          description:"Optimistic version of address row"` // Optimistic version of address row
	ReplacedFromAddressId uint64      `json:"replacedFromAddressId" orm:"replaced_from_address_id" description:"Address replaced chain source"`     // Address replaced chain source
	Label                 string      `json:"label"                 orm:"label"                    description:"Address label"`                     // Address label
	ReceiverName          string      `json:"receiverName"          orm:"receiver_name"            description:"Receiver name"`                     // Receiver name
	ReceiverPhone         string      `json:"receiverPhone"         orm:"receiver_phone"           description:"Receiver phone"`                    // Receiver phone
	CountryCode           string      `json:"countryCode"           orm:"country_code"             description:"ISO country code"`                  // ISO country code
	ProvinceCode          string      `json:"provinceCode"          orm:"province_code"            description:"Province code"`                     // Province code
	ProvinceName          string      `json:"provinceName"          orm:"province_name"            description:"Province name"`                     // Province name
	CityCode              string      `json:"cityCode"              orm:"city_code"                description:"City code"`                         // City code
	CityName              string      `json:"cityName"              orm:"city_name"                description:"City name"`                         // City name
	DistrictCode          string      `json:"districtCode"          orm:"district_code"            description:"District code"`                     // District code
	DistrictName          string      `json:"districtName"          orm:"district_name"            description:"District name"`                     // District name
	Street                string      `json:"street"                orm:"street"                   description:"Street/town"`                       // Street/town
	Detail                string      `json:"detail"                orm:"detail"                   description:"Detailed address"`                  // Detailed address
	PostalCode            string      `json:"postalCode"            orm:"postal_code"              description:"Postal code"`                       // Postal code
	IsDefault             int         `json:"isDefault"             orm:"is_default"               description:"1 default, NULL non-default"`       // 1 default, NULL non-default
	DefaultSlot           int         `json:"defaultSlot"           orm:"default_slot"             description:"Default uniqueness slot"`           // Default uniqueness slot
	Latitude              float64     `json:"latitude"              orm:"latitude"                 description:"Latitude"`                          // Latitude
	Longitude             float64     `json:"longitude"             orm:"longitude"                description:"Longitude"`                         // Longitude
	Ext                   string      `json:"ext"                   orm:"ext"                      description:"Extension payload"`                 // Extension payload
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"               description:""`                                  //
	UpdatedAt             *gtime.Time `json:"updatedAt"             orm:"updated_at"               description:""`                                  //
}
