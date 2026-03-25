// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerEntityDoc is the golang structure for table seller_entity_doc.
type SellerEntityDoc struct {
	Id          uint64      `json:"id"          orm:"id"            ` //
	EntityId    uint64      `json:"entityId"    orm:"entity_id"     ` // FK seller_entity.id
	EntityNo    string      `json:"entityNo"    orm:"entity_no"     ` // Redundant external entity no
	DocTypeCode string      `json:"docTypeCode" orm:"doc_type_code" ` // BUSINESS_LICENSE/ID_FRONT/...
	AssetId     uint64      `json:"assetId"     orm:"asset_id"      ` // media-svc asset id
	ValidFrom   *gtime.Time `json:"validFrom"   orm:"valid_from"    ` //
	ValidUntil  *gtime.Time `json:"validUntil"  orm:"valid_until"   ` //
	Issuer      string      `json:"issuer"      orm:"issuer"        ` //
	ExtJson     string      `json:"extJson"     orm:"ext_json"      ` //
	Status      uint        `json:"status"      orm:"status"        ` // 1 ACTIVE,2 INACTIVE
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    ` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    ` //
}
