// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerEntity is the golang structure of table seller_entity for DAO operations like Where/Data.
type SellerEntity struct {
	g.Meta                    `orm:"table:seller_entity, do:true"`
	Id                        any         // Internal PK
	EntityNo                  any         // External business id
	OwnerUserId               any         // IAM user id
	MerchantTypeCode          any         // PERSONAL/ENTERPRISE/...
	EntityName                any         // Display name
	ContactName               any         // Contact person
	ContactPhone              any         // Contact phone
	ContactEmail              any         // Contact email
	SubjectTypeCode           any         // Legal subject type
	SubjectName               any         // Legal subject name
	SubjectCertNo             any         // Legal certificate no
	LegalRepresentativeName   any         // Legal representative
	LegalRepresentativeCertNo any         // Legal rep certificate no
	CertValidFrom             *gtime.Time // Legal cert valid from
	CertValidUntil            *gtime.Time // Legal cert valid until
	LegalSubjectExtJson       any         // Ext fields for legal subject
	ExtJson                   any         // Entity ext fields
	Version                   any         // Optimistic version
	CreatedAt                 *gtime.Time //
	UpdatedAt                 *gtime.Time //
}
