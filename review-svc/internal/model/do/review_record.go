// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReviewRecord is the golang structure of table review_record for DAO operations like Where/Data.
type ReviewRecord struct {
	g.Meta           `orm:"table:review_record, do:true"`
	Id               any         //
	ReviewNo         any         //
	OrderNo          any         //
	SubOrderNo       any         //
	ItemNo           any         //
	UserId           any         //
	ShopNo           any         //
	SpuNo            any         //
	SkuNo            any         //
	Score            any         //
	Content          any         //
	MediasJson       any         //
	Anonymous        any         //
	AppendContent    any         //
	AppendMediasJson any         //
	AppendAt         *gtime.Time //
	SellerReply      any         //
	SellerReplyAt    *gtime.Time //
	ReviewStatus     any         //
	LikeCount        any         //
	Version          any         //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
	DeletedAt        *gtime.Time //
}
