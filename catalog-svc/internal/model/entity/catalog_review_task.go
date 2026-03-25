// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogReviewTask is the golang structure for table catalog_review_task.
type CatalogReviewTask struct {
	Id                 uint64      `json:"id"                 orm:"id"                    ` //
	TaskNo             string      `json:"taskNo"             orm:"task_no"               ` //
	SpuNo              string      `json:"spuNo"              orm:"spu_no"                ` //
	ShopNo             string      `json:"shopNo"             orm:"shop_no"               ` //
	SpuVersionAtSubmit uint        `json:"spuVersionAtSubmit" orm:"spu_version_at_submit" ` //
	SubmitNote         string      `json:"submitNote"         orm:"submit_note"           ` //
	ReviewStatus       uint        `json:"reviewStatus"       orm:"review_status"         ` //
	RejectReasonCode   string      `json:"rejectReasonCode"   orm:"reject_reason_code"    ` //
	RejectComment      string      `json:"rejectComment"      orm:"reject_comment"        ` //
	ReviewComment      string      `json:"reviewComment"      orm:"review_comment"        ` //
	ReviewerId         uint64      `json:"reviewerId"         orm:"reviewer_id"           ` //
	SubmittedAt        *gtime.Time `json:"submittedAt"        orm:"submitted_at"          ` //
	ReviewStartedAt    *gtime.Time `json:"reviewStartedAt"    orm:"review_started_at"     ` //
	ReviewedAt         *gtime.Time `json:"reviewedAt"         orm:"reviewed_at"           ` //
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"            ` //
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"            ` //
}
