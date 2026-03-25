// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SellerShopDao is the data access object for the table seller_shop.
type SellerShopDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SellerShopColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SellerShopColumns defines and stores column names for the table seller_shop.
type SellerShopColumns struct {
	Id                  string // Internal shop id (scope_id)
	ShopNo              string // External business shop id
	EntityId            string // FK seller_entity.id
	EntityNo            string // External entity id
	ApplicationNo       string // Source application no
	OwnerUserId         string // IAM user id
	ShopName            string // Requested unique shop name
	ShopNameNorm        string // Normalized lower name
	ShopDisplayName     string // Display name
	ShopTypeCode        string // Shop type code
	MainCategoryIdsJson string // Main category id list
	LogoAssetId         string //
	BannerAssetId       string //
	ServicePhone        string //
	ServiceEmail        string //
	Description         string //
	ExtJson             string //
	Status              string // ShopStatus enum
	BuyerVisible        string // 1 visible to buyers
	Version             string // Optimistic version
	FreezeReasonCode    string //
	FreezeReason        string //
	CloseReasonCode     string //
	CloseReason         string //
	ProvisionStartedAt  string //
	ProvisionFinishedAt string //
	ActivatedAt         string //
	FrozenAt            string //
	ClosedAt            string //
	CreatedAt           string //
	UpdatedAt           string //
}

// sellerShopColumns holds the columns for the table seller_shop.
var sellerShopColumns = SellerShopColumns{
	Id:                  "id",
	ShopNo:              "shop_no",
	EntityId:            "entity_id",
	EntityNo:            "entity_no",
	ApplicationNo:       "application_no",
	OwnerUserId:         "owner_user_id",
	ShopName:            "shop_name",
	ShopNameNorm:        "shop_name_norm",
	ShopDisplayName:     "shop_display_name",
	ShopTypeCode:        "shop_type_code",
	MainCategoryIdsJson: "main_category_ids_json",
	LogoAssetId:         "logo_asset_id",
	BannerAssetId:       "banner_asset_id",
	ServicePhone:        "service_phone",
	ServiceEmail:        "service_email",
	Description:         "description",
	ExtJson:             "ext_json",
	Status:              "status",
	BuyerVisible:        "buyer_visible",
	Version:             "version",
	FreezeReasonCode:    "freeze_reason_code",
	FreezeReason:        "freeze_reason",
	CloseReasonCode:     "close_reason_code",
	CloseReason:         "close_reason",
	ProvisionStartedAt:  "provision_started_at",
	ProvisionFinishedAt: "provision_finished_at",
	ActivatedAt:         "activated_at",
	FrozenAt:            "frozen_at",
	ClosedAt:            "closed_at",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewSellerShopDao creates and returns a new DAO object for table data access.
func NewSellerShopDao(handlers ...gdb.ModelHandler) *SellerShopDao {
	return &SellerShopDao{
		group:    "default",
		table:    "seller_shop",
		columns:  sellerShopColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SellerShopDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SellerShopDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SellerShopDao) Columns() SellerShopColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SellerShopDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SellerShopDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SellerShopDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
