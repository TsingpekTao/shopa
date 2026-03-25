// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SellerEntityDao is the data access object for the table seller_entity.
type SellerEntityDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  SellerEntityColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// SellerEntityColumns defines and stores column names for the table seller_entity.
type SellerEntityColumns struct {
	Id                        string // Internal PK
	EntityNo                  string // External business id
	OwnerUserId               string // IAM user id
	MerchantTypeCode          string // PERSONAL/ENTERPRISE/...
	EntityName                string // Display name
	ContactName               string // Contact person
	ContactPhone              string // Contact phone
	ContactEmail              string // Contact email
	SubjectTypeCode           string // Legal subject type
	SubjectName               string // Legal subject name
	SubjectCertNo             string // Legal certificate no
	LegalRepresentativeName   string // Legal representative
	LegalRepresentativeCertNo string // Legal rep certificate no
	CertValidFrom             string // Legal cert valid from
	CertValidUntil            string // Legal cert valid until
	LegalSubjectExtJson       string // Ext fields for legal subject
	ExtJson                   string // Entity ext fields
	Version                   string // Optimistic version
	CreatedAt                 string //
	UpdatedAt                 string //
}

// sellerEntityColumns holds the columns for the table seller_entity.
var sellerEntityColumns = SellerEntityColumns{
	Id:                        "id",
	EntityNo:                  "entity_no",
	OwnerUserId:               "owner_user_id",
	MerchantTypeCode:          "merchant_type_code",
	EntityName:                "entity_name",
	ContactName:               "contact_name",
	ContactPhone:              "contact_phone",
	ContactEmail:              "contact_email",
	SubjectTypeCode:           "subject_type_code",
	SubjectName:               "subject_name",
	SubjectCertNo:             "subject_cert_no",
	LegalRepresentativeName:   "legal_representative_name",
	LegalRepresentativeCertNo: "legal_representative_cert_no",
	CertValidFrom:             "cert_valid_from",
	CertValidUntil:            "cert_valid_until",
	LegalSubjectExtJson:       "legal_subject_ext_json",
	ExtJson:                   "ext_json",
	Version:                   "version",
	CreatedAt:                 "created_at",
	UpdatedAt:                 "updated_at",
}

// NewSellerEntityDao creates and returns a new DAO object for table data access.
func NewSellerEntityDao(handlers ...gdb.ModelHandler) *SellerEntityDao {
	return &SellerEntityDao{
		group:    "default",
		table:    "seller_entity",
		columns:  sellerEntityColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SellerEntityDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SellerEntityDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SellerEntityDao) Columns() SellerEntityColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SellerEntityDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SellerEntityDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SellerEntityDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
