// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsConsumptionDetail is the golang structure of table points_consumption_detail for DAO operations like Where/Data.
type PointsConsumptionDetail struct {
	g.Meta         `orm:"table:points_consumption_detail, do:true"`
	Id             any         //
	DetailNo       any         //
	UserId         any         //
	ReservationNo  any         //
	OrderNo        any         //
	RefundNo       any         //
	SubOrderNo     any         //
	BucketNo       any         //
	DetailStatusCode any       // LOCKED/CONFIRMED/CANCELED
	ConsumedPoints any         //
	ReturnedPoints any         //
	CashAmountCent any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
