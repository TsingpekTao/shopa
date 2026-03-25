// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogReviewTask is the golang structure of table catalog_review_task for DAO operations like Where/Data.
type CatalogReviewTask struct {
	g.Meta             `orm:"table:catalog_review_task, do:true"`
	Id                 any         //
	TaskNo             any         //
	SpuNo              any         //
	ShopNo             any         //
	SpuVersionAtSubmit any         //
	SubmitNote         any         //
	ReviewStatus       any         //
	RejectReasonCode   any         //
	RejectComment      any         //
	ReviewComment      any         //
	ReviewerId         any         //
	SubmittedAt        *gtime.Time //
	ReviewStartedAt    *gtime.Time //
	ReviewedAt         *gtime.Time //
	CreatedAt          *gtime.Time //
	UpdatedAt          *gtime.Time //
}
