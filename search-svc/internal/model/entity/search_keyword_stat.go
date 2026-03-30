// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SearchKeywordStat is the golang structure for table search_keyword_stat.
type SearchKeywordStat struct {
	Id             uint64      `json:"id"             orm:"id"               ` //
	Keyword        string      `json:"keyword"        orm:"keyword"          ` //
	SearchCount    uint64      `json:"searchCount"    orm:"search_count"     ` //
	LastSearchedAt *gtime.Time `json:"lastSearchedAt" orm:"last_searched_at" ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       ` //
}
