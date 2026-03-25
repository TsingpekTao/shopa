// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerEntityDoc is the golang structure of table seller_entity_doc for DAO operations like Where/Data.
type SellerEntityDoc struct {
	g.Meta      `orm:"table:seller_entity_doc, do:true"`
	Id          any         //
	EntityId    any         // FK seller_entity.id
	EntityNo    any         // Redundant external entity no
	DocTypeCode any         // BUSINESS_LICENSE/ID_FRONT/...
	AssetId     any         // media-svc asset id
	ValidFrom   *gtime.Time //
	ValidUntil  *gtime.Time //
	Issuer      any         //
	ExtJson     any         //
	Status      any         // 1 ACTIVE,2 INACTIVE
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
