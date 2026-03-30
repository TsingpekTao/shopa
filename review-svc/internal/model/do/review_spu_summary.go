// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewSpuSummary is the golang structure of table review_spu_summary for DAO operations like Where/Data.
type ReviewSpuSummary struct {
	g.Meta           `orm:"table:review_spu_summary, do:true"`
	Id               any         //
	SpuNo            any         //
	TotalReviews     any         //
	Score1Count      any         //
	Score2Count      any         //
	Score3Count      any         //
	Score4Count      any         //
	Score5Count      any         //
	AvgScoreX100     any         //
	PositiveRateX100 any         //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
	DeletedAt        *gtime.Time //
}
