// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewSpuSummary is the golang structure for table review_spu_summary.
type ReviewSpuSummary struct {
	Id               uint64      `json:"id"               orm:"id"                 description:""` //
	SpuNo            string      `json:"spuNo"            orm:"spu_no"             description:""` //
	TotalReviews     uint64      `json:"totalReviews"     orm:"total_reviews"      description:""` //
	Score1Count      uint64      `json:"score1Count"      orm:"score_1_count"      description:""` //
	Score2Count      uint64      `json:"score2Count"      orm:"score_2_count"      description:""` //
	Score3Count      uint64      `json:"score3Count"      orm:"score_3_count"      description:""` //
	Score4Count      uint64      `json:"score4Count"      orm:"score_4_count"      description:""` //
	Score5Count      uint64      `json:"score5Count"      orm:"score_5_count"      description:""` //
	AvgScoreX100     uint64      `json:"avgScoreX100"     orm:"avg_score_x100"     description:""` //
	PositiveRateX100 uint64      `json:"positiveRateX100" orm:"positive_rate_x100" description:""` //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:""` //
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         description:""` //
	DeletedAt        *gtime.Time `json:"deletedAt"        orm:"deleted_at"         description:""` //
}
