// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SellerEntityDocDao is the data access object for the table seller_entity_doc.
type SellerEntityDocDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  SellerEntityDocColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// SellerEntityDocColumns defines and stores column names for the table seller_entity_doc.
type SellerEntityDocColumns struct {
	Id          string //
	EntityId    string // FK seller_entity.id
	EntityNo    string // Redundant external entity no
	DocTypeCode string // BUSINESS_LICENSE/ID_FRONT/...
	AssetId     string // media-svc asset id
	ValidFrom   string //
	ValidUntil  string //
	Issuer      string //
	ExtJson     string //
	Status      string // 1 ACTIVE,2 INACTIVE
	CreatedAt   string //
	UpdatedAt   string //
}

// sellerEntityDocColumns holds the columns for the table seller_entity_doc.
var sellerEntityDocColumns = SellerEntityDocColumns{
	Id:          "id",
	EntityId:    "entity_id",
	EntityNo:    "entity_no",
	DocTypeCode: "doc_type_code",
	AssetId:     "asset_id",
	ValidFrom:   "valid_from",
	ValidUntil:  "valid_until",
	Issuer:      "issuer",
	ExtJson:     "ext_json",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewSellerEntityDocDao creates and returns a new DAO object for table data access.
func NewSellerEntityDocDao(handlers ...gdb.ModelHandler) *SellerEntityDocDao {
	return &SellerEntityDocDao{
		group:    "default",
		table:    "seller_entity_doc",
		columns:  sellerEntityDocColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SellerEntityDocDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SellerEntityDocDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SellerEntityDocDao) Columns() SellerEntityDocColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SellerEntityDocDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SellerEntityDocDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SellerEntityDocDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
