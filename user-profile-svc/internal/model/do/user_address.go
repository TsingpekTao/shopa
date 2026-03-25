// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAddress is the golang structure of table user_address for DAO operations like Where/Data.
type UserAddress struct {
	g.Meta                `orm:"table:user_address, do:true"`
	AddressId             any         // Address ID
	UserId                any         // Owner user ID
	Status                any         // 1=active,2=deleted,3=replaced
	AddressVersion        any         // Optimistic version of address row
	ReplacedFromAddressId any         // Address replaced chain source
	Label                 any         // Address label
	ReceiverName          any         // Receiver name
	ReceiverPhone         any         // Receiver phone
	CountryCode           any         // ISO country code
	ProvinceCode          any         // Province code
	ProvinceName          any         // Province name
	CityCode              any         // City code
	CityName              any         // City name
	DistrictCode          any         // District code
	DistrictName          any         // District name
	Street                any         // Street/town
	Detail                any         // Detailed address
	PostalCode            any         // Postal code
	IsDefault             any         // 1 default, NULL non-default
	DefaultSlot           any         // Default uniqueness slot
	Latitude              any         // Latitude
	Longitude             any         // Longitude
	Ext                   any         // Extension payload
	CreatedAt             *gtime.Time //
	UpdatedAt             *gtime.Time //
}
