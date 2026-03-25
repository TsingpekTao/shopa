// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CatalogBrand is the golang structure for table catalog_brand.
type CatalogBrand struct {
	Id        uint64      `json:"id"        orm:"id"         ` //
	BrandNo   string      `json:"brandNo"   orm:"brand_no"   ` //
	BrandName string      `json:"brandName" orm:"brand_name" ` //
	Status    uint        `json:"status"    orm:"status"     ` // 1 ENABLED, 2 DISABLED
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
