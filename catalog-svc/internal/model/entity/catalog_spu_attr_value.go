// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogSpuAttrValue is the golang structure for table catalog_spu_attr_value.
type CatalogSpuAttrValue struct {
	Id        uint64      `json:"id"        orm:"id"         ` //
	SpuNo     string      `json:"spuNo"     orm:"spu_no"     ` //
	AttrCode  string      `json:"attrCode"  orm:"attr_code"  ` //
	AttrName  string      `json:"attrName"  orm:"attr_name"  ` //
	AttrScope uint        `json:"attrScope" orm:"attr_scope" ` //
	AttrValue string      `json:"attrValue" orm:"attr_value" ` //
	SortOrder int         `json:"sortOrder" orm:"sort_order" ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
