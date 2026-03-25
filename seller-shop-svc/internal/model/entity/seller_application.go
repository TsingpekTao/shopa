// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerApplication is the golang structure for table seller_application.
type SellerApplication struct {
	Id                    uint64      `json:"id"                    orm:"id"                      ` //
	ApplicationNo         string      `json:"applicationNo"         orm:"application_no"          ` // External application id
	Version               uint        `json:"version"               orm:"version"                 ` // Optimistic version
	PreviousApplicationNo string      `json:"previousApplicationNo" orm:"previous_application_no" ` // Chain link to previous rejected application
	NextApplicationNo     string      `json:"nextApplicationNo"     orm:"next_application_no"     ` // Chain link to next resubmitted application
	OwnerUserId           uint64      `json:"ownerUserId"           orm:"owner_user_id"           ` //
	Status                uint        `json:"status"                orm:"status"                  ` // ApplicationStatus enum
	EntityNo              string      `json:"entityNo"              orm:"entity_no"               ` //
	ShopNo                string      `json:"shopNo"                orm:"shop_no"                 ` //
	EntityName            string      `json:"entityName"            orm:"entity_name"             ` // Denormalized for list/search
	ShopName              string      `json:"shopName"              orm:"shop_name"               ` // Denormalized for list/search
	ShopNameNorm          string      `json:"shopNameNorm"          orm:"shop_name_norm"          ` // Denormalized normalized shop name
	EntityDraftJson       string      `json:"entityDraftJson"       orm:"entity_draft_json"       ` // Current draft entity payload
	ShopDraftJson         string      `json:"shopDraftJson"         orm:"shop_draft_json"         ` // Current draft shop payload
	EntitySubmittedJson   string      `json:"entitySubmittedJson"   orm:"entity_submitted_json"   ` // Snapshot for review at submit time
	ShopSubmittedJson     string      `json:"shopSubmittedJson"     orm:"shop_submitted_json"     ` // Snapshot for review at submit time
	LatestRejectJson      string      `json:"latestRejectJson"      orm:"latest_reject_json"      ` // Latest reject info
	SubmittedAt           *gtime.Time `json:"submittedAt"           orm:"submitted_at"            ` //
	ReviewStartedAt       *gtime.Time `json:"reviewStartedAt"       orm:"review_started_at"       ` //
	ReviewedAt            *gtime.Time `json:"reviewedAt"            orm:"reviewed_at"             ` //
	ReviewerId            uint64      `json:"reviewerId"            orm:"reviewer_id"             ` //
	ReviewComment         string      `json:"reviewComment"         orm:"review_comment"          ` //
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"              ` //
	UpdatedAt             *gtime.Time `json:"updatedAt"             orm:"updated_at"              ` //
}
