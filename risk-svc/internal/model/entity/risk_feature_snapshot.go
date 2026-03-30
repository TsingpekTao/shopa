// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskFeatureSnapshot is the golang structure for table risk_feature_snapshot.
type RiskFeatureSnapshot struct {
	Id           uint64      `json:"id"           orm:"id"            ` //
	UserId       uint64      `json:"userId"       orm:"user_id"       ` //
	RiskScore    uint        `json:"riskScore"    orm:"risk_score"    ` //
	TagsJson     string      `json:"tagsJson"     orm:"tags_json"     ` //
	FeaturesJson string      `json:"featuresJson" orm:"features_json" ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
}
