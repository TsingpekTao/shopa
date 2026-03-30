// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentTrackingNode is the golang structure for table fulfillment_tracking_node.
type FulfillmentTrackingNode struct {
	Id         uint64      `json:"id"         orm:"id"          description:""` //
	NodeNo     string      `json:"nodeNo"     orm:"node_no"     description:""` //
	ShipmentNo string      `json:"shipmentNo" orm:"shipment_no" description:""` //
	StatusCode string      `json:"statusCode" orm:"status_code" description:""` //
	Content    string      `json:"content"    orm:"content"     description:""` //
	Location   string      `json:"location"   orm:"location"    description:""` //
	EventTime  *gtime.Time `json:"eventTime"  orm:"event_time"  description:""` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""` //
}
