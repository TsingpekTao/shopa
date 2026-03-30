// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderAddressSnapshot is the golang structure of table order_address_snapshot for DAO operations like Where/Data.
type OrderAddressSnapshot struct {
	g.Meta               `orm:"table:order_address_snapshot, do:true"`
	Id                   any         //
	OrderNo              any         //
	SourceAddressId      any         //
	SourceAddressVersion any         //
	ReceiverName         any         //
	ReceiverPhone        any         //
	CountryCode          any         //
	ProvinceCode         any         //
	ProvinceName         any         //
	CityCode             any         //
	CityName             any         //
	DistrictCode         any         //
	DistrictName         any         //
	Street               any         //
	Detail               any         //
	PostalCode           any         //
	Latitude             any         //
	Longitude            any         //
	CreatedAt            *gtime.Time //
	UpdatedAt            *gtime.Time //
}
