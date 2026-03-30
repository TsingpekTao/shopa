// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskFeatureSnapshot is the golang structure of table risk_feature_snapshot for DAO operations like Where/Data.
type RiskFeatureSnapshot struct {
	g.Meta       `orm:"table:risk_feature_snapshot, do:true"`
	Id           any         //
	UserId       any         //
	RiskScore    any         //
	TagsJson     any         //
	FeaturesJson any         //
	UpdatedAt    *gtime.Time //
}
