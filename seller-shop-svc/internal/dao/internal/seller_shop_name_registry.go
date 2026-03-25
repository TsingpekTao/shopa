// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SellerShopNameRegistryDao is the data access object for the table seller_shop_name_registry.
type SellerShopNameRegistryDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  SellerShopNameRegistryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// SellerShopNameRegistryColumns defines and stores column names for the table seller_shop_name_registry.
type SellerShopNameRegistryColumns struct {
	Id            string //
	ShopNameNorm  string // Normalized unique shop name
	ShopName      string // Original input shop name
	OwnerUserId   string //
	ApplicationNo string //
	Status        string //
	ExpireAt      string // Reservation expiration time
	BoundShopNo   string // Final bound shop_no when approved
	CreatedAt     string //
	UpdatedAt     string //
}

// sellerShopNameRegistryColumns holds the columns for the table seller_shop_name_registry.
var sellerShopNameRegistryColumns = SellerShopNameRegistryColumns{
	Id:            "id",
	ShopNameNorm:  "shop_name_norm",
	ShopName:      "shop_name",
	OwnerUserId:   "owner_user_id",
	ApplicationNo: "application_no",
	Status:        "status",
	ExpireAt:      "expire_at",
	BoundShopNo:   "bound_shop_no",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewSellerShopNameRegistryDao creates and returns a new DAO object for table data access.
func NewSellerShopNameRegistryDao(handlers ...gdb.ModelHandler) *SellerShopNameRegistryDao {
	return &SellerShopNameRegistryDao{
		group:    "default",
		table:    "seller_shop_name_registry",
		columns:  sellerShopNameRegistryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SellerShopNameRegistryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SellerShopNameRegistryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SellerShopNameRegistryDao) Columns() SellerShopNameRegistryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SellerShopNameRegistryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SellerShopNameRegistryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SellerShopNameRegistryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
