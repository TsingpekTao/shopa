// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsConsumptionDetail is the golang structure for table points_consumption_detail.
type PointsConsumptionDetail struct {
	Id             uint64      `json:"id"             orm:"id"               ` //
	DetailNo       string      `json:"detailNo"       orm:"detail_no"        ` //
	UserId         uint64      `json:"userId"         orm:"user_id"          ` //
	ReservationNo  string      `json:"reservationNo"  orm:"reservation_no"   ` //
	OrderNo        string      `json:"orderNo"        orm:"order_no"         ` //
	RefundNo       string      `json:"refundNo"       orm:"refund_no"        ` //
	SubOrderNo     string      `json:"subOrderNo"     orm:"sub_order_no"     ` //
	BucketNo       string      `json:"bucketNo"       orm:"bucket_no"        ` //
	DetailStatusCode string    `json:"detailStatusCode" orm:"detail_status_code" ` // LOCKED/CONFIRMED/CANCELED
	ConsumedPoints uint64      `json:"consumedPoints" orm:"consumed_points"  ` //
	ReturnedPoints uint64      `json:"returnedPoints" orm:"returned_points"  ` //
	CashAmountCent int64       `json:"cashAmountCent" orm:"cash_amount_cent" ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       ` //
}
