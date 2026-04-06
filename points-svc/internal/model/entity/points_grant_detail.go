// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsGrantDetail is the golang structure for table points_grant_detail.
type PointsGrantDetail struct {
	Id              uint64      `json:"id"              orm:"id"                ` //
	GrantDetailNo   string      `json:"grantDetailNo"   orm:"grant_detail_no"   ` //
	UserId          uint64      `json:"userId"          orm:"user_id"           ` //
	OrderNo         string      `json:"orderNo"         orm:"order_no"          ` //
	SubOrderNo      string      `json:"subOrderNo"      orm:"sub_order_no"      ` //
	ShopNo          string      `json:"shopNo"          orm:"shop_no"           ` //
	GrantedPoints   uint64      `json:"grantedPoints"   orm:"granted_points"    ` //
	ReversedPoints  uint64      `json:"reversedPoints"  orm:"reversed_points"   ` //
	RuleCode        string      `json:"ruleCode"        orm:"rule_code"         ` //
	RuleSnapshotJson string     `json:"ruleSnapshotJson" orm:"rule_snapshot_json" ` //
	GrantStatusCode string      `json:"grantStatusCode" orm:"grant_status_code" ` // GRANTED/PARTIAL_REVERSED/REVERSED
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        ` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"        ` //
}
