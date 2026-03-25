// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CatalogAttributeTemplateDao is the data access object for the table catalog_attribute_template.
type CatalogAttributeTemplateDao struct {
	table    string                          // table is the underlying table name of the DAO.
	group    string                          // group is the database configuration group name of the current DAO.
	columns  CatalogAttributeTemplateColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler              // handlers for customized model modification.
}

// CatalogAttributeTemplateColumns defines and stores column names for the table catalog_attribute_template.
type CatalogAttributeTemplateColumns struct {
	Id             string //
	CategoryId     string //
	AttrCode       string //
	AttrName       string //
	AttrScope      string // 1 SPU, 2 SKU_SALE, 3 EXT
	ValueType      string // TEXT/NUMBER/BOOL/ENUM
	RequiredFlag   string //
	SearchableFlag string //
	FilterableFlag string //
	OptionsJson    string // Enum options if needed
	SortOrder      string //
	Status         string // 1 ENABLED, 2 DISABLED
	CreatedAt      string //
	UpdatedAt      string //
}

// catalogAttributeTemplateColumns holds the columns for the table catalog_attribute_template.
var catalogAttributeTemplateColumns = CatalogAttributeTemplateColumns{
	Id:             "id",
	CategoryId:     "category_id",
	AttrCode:       "attr_code",
	AttrName:       "attr_name",
	AttrScope:      "attr_scope",
	ValueType:      "value_type",
	RequiredFlag:   "required_flag",
	SearchableFlag: "searchable_flag",
	FilterableFlag: "filterable_flag",
	OptionsJson:    "options_json",
	SortOrder:      "sort_order",
	Status:         "status",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewCatalogAttributeTemplateDao creates and returns a new DAO object for table data access.
func NewCatalogAttributeTemplateDao(handlers ...gdb.ModelHandler) *CatalogAttributeTemplateDao {
	return &CatalogAttributeTemplateDao{
		group:    "default",
		table:    "catalog_attribute_template",
		columns:  catalogAttributeTemplateColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CatalogAttributeTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CatalogAttributeTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CatalogAttributeTemplateDao) Columns() CatalogAttributeTemplateColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CatalogAttributeTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CatalogAttributeTemplateDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CatalogAttributeTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
