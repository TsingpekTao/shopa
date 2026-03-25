// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogAttributeTemplate is the golang structure for table catalog_attribute_template.
type CatalogAttributeTemplate struct {
	Id             uint64      `json:"id"             orm:"id"              ` //
	CategoryId     uint64      `json:"categoryId"     orm:"category_id"     ` //
	AttrCode       string      `json:"attrCode"       orm:"attr_code"       ` //
	AttrName       string      `json:"attrName"       orm:"attr_name"       ` //
	AttrScope      uint        `json:"attrScope"      orm:"attr_scope"      ` // 1 SPU, 2 SKU_SALE, 3 EXT
	ValueType      string      `json:"valueType"      orm:"value_type"      ` // TEXT/NUMBER/BOOL/ENUM
	RequiredFlag   int         `json:"requiredFlag"   orm:"required_flag"   ` //
	SearchableFlag int         `json:"searchableFlag" orm:"searchable_flag" ` //
	FilterableFlag int         `json:"filterableFlag" orm:"filterable_flag" ` //
	OptionsJson    string      `json:"optionsJson"    orm:"options_json"    ` // Enum options if needed
	SortOrder      int         `json:"sortOrder"      orm:"sort_order"      ` //
	Status         uint        `json:"status"         orm:"status"          ` // 1 ENABLED, 2 DISABLED
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
}
