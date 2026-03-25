// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogIdempotencyDao is the data access object for the table catalog_idempotency.
type CatalogIdempotencyDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  CatalogIdempotencyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// CatalogIdempotencyColumns defines and stores column names for the table catalog_idempotency.
type CatalogIdempotencyColumns struct {
	Id             string //
	ActorUserId    string //
	ActionCode     string //
	IdempotencyKey string //
	RequestHash    string //
	ShopNo         string //
	SpuNo          string //
	Status         string //
	ResponseCode   string //
	ResponseJson   string //
	ExpiredAt      string //
	CreatedAt      string //
	UpdatedAt      string //
}

// catalogIdempotencyColumns holds the columns for the table catalog_idempotency.
var catalogIdempotencyColumns = CatalogIdempotencyColumns{
	Id:             "id",
	ActorUserId:    "actor_user_id",
	ActionCode:     "action_code",
	IdempotencyKey: "idempotency_key",
	RequestHash:    "request_hash",
	ShopNo:         "shop_no",
	SpuNo:          "spu_no",
	Status:         "status",
	ResponseCode:   "response_code",
	ResponseJson:   "response_json",
	ExpiredAt:      "expired_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewCatalogIdempotencyDao creates and returns a new DAO object for table data access.
func NewCatalogIdempotencyDao(handlers ...gdb.ModelHandler) *CatalogIdempotencyDao {
	return &CatalogIdempotencyDao{
		group:    "default",
		table:    "catalog_idempotency",
		columns:  catalogIdempotencyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogIdempotencyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogIdempotencyDao) Columns() CatalogIdempotencyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogIdempotencyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
