// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FulfillmentTrackingNode is the golang structure of table fulfillment_tracking_node for DAO operations like Where/Data.
type FulfillmentTrackingNode struct {
	g.Meta     `orm:"table:fulfillment_tracking_node, do:true"`
	Id         any         //
	NodeNo     any         //
	ShipmentNo any         //
	StatusCode any         //
	Content    any         //
	Location   any         //
	EventTime  *gtime.Time //
	CreatedAt  *gtime.Time //
}
