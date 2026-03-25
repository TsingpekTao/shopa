// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SellerEntity is the golang structure for table seller_entity.
type SellerEntity struct {
	Id                        uint64      `json:"id"                        orm:"id"                           ` // Internal PK
	EntityNo                  string      `json:"entityNo"                  orm:"entity_no"                    ` // External business id
	OwnerUserId               uint64      `json:"ownerUserId"               orm:"owner_user_id"                ` // IAM user id
	MerchantTypeCode          string      `json:"merchantTypeCode"          orm:"merchant_type_code"           ` // PERSONAL/ENTERPRISE/...
	EntityName                string      `json:"entityName"                orm:"entity_name"                  ` // Display name
	ContactName               string      `json:"contactName"               orm:"contact_name"                 ` // Contact person
	ContactPhone              string      `json:"contactPhone"              orm:"contact_phone"                ` // Contact phone
	ContactEmail              string      `json:"contactEmail"              orm:"contact_email"                ` // Contact email
	SubjectTypeCode           string      `json:"subjectTypeCode"           orm:"subject_type_code"            ` // Legal subject type
	SubjectName               string      `json:"subjectName"               orm:"subject_name"                 ` // Legal subject name
	SubjectCertNo             string      `json:"subjectCertNo"             orm:"subject_cert_no"              ` // Legal certificate no
	LegalRepresentativeName   string      `json:"legalRepresentativeName"   orm:"legal_representative_name"    ` // Legal representative
	LegalRepresentativeCertNo string      `json:"legalRepresentativeCertNo" orm:"legal_representative_cert_no" ` // Legal rep certificate no
	CertValidFrom             *gtime.Time `json:"certValidFrom"             orm:"cert_valid_from"              ` // Legal cert valid from
	CertValidUntil            *gtime.Time `json:"certValidUntil"            orm:"cert_valid_until"             ` // Legal cert valid until
	LegalSubjectExtJson       string      `json:"legalSubjectExtJson"       orm:"legal_subject_ext_json"       ` // Ext fields for legal subject
	ExtJson                   string      `json:"extJson"                   orm:"ext_json"                     ` // Entity ext fields
	Version                   uint        `json:"version"                   orm:"version"                      ` // Optimistic version
	CreatedAt                 *gtime.Time `json:"createdAt"                 orm:"created_at"                   ` //
	UpdatedAt                 *gtime.Time `json:"updatedAt"                 orm:"updated_at"                   ` //
}
