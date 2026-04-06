// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsGrantDetail is the golang structure of table points_grant_detail for DAO operations like Where/Data.
type PointsGrantDetail struct {
	g.Meta          `orm:"table:points_grant_detail, do:true"`
	Id              any         //
	GrantDetailNo   any         //
	UserId          any         //
	OrderNo         any         //
	SubOrderNo      any         //
	ShopNo          any         //
	GrantedPoints   any         //
	ReversedPoints  any         //
	RuleCode        any         //
	RuleSnapshotJson any        //
	GrantStatusCode any         // GRANTED/PARTIAL_REVERSED/REVERSED
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
