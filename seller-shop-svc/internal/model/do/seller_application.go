// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerApplication is the golang structure of table seller_application for DAO operations like Where/Data.
type SellerApplication struct {
	g.Meta                `orm:"table:seller_application, do:true"`
	Id                    any         //
	ApplicationNo         any         // External application id
	Version               any         // Optimistic version
	PreviousApplicationNo any         // Chain link to previous rejected application
	NextApplicationNo     any         // Chain link to next resubmitted application
	OwnerUserId           any         //
	Status                any         // ApplicationStatus enum
	EntityNo              any         //
	ShopNo                any         //
	EntityName            any         // Denormalized for list/search
	ShopName              any         // Denormalized for list/search
	ShopNameNorm          any         // Denormalized normalized shop name
	EntityDraftJson       any         // Current draft entity payload
	ShopDraftJson         any         // Current draft shop payload
	EntitySubmittedJson   any         // Snapshot for review at submit time
	ShopSubmittedJson     any         // Snapshot for review at submit time
	LatestRejectJson      any         // Latest reject info
	SubmittedAt           *gtime.Time //
	ReviewStartedAt       *gtime.Time //
	ReviewedAt            *gtime.Time //
	ReviewerId            any         //
	ReviewComment         any         //
	CreatedAt             *gtime.Time //
	UpdatedAt             *gtime.Time //
}
