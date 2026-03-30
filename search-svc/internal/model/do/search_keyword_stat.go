// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SearchKeywordStat is the golang structure of table search_keyword_stat for DAO operations like Where/Data.
type SearchKeywordStat struct {
	g.Meta         `orm:"table:search_keyword_stat, do:true"`
	Id             any         //
	Keyword        any         //
	SearchCount    any         //
	LastSearchedAt *gtime.Time //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
