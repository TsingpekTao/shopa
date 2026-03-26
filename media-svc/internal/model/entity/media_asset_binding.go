// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaAssetBinding 是 media_asset_binding 表的结构体。
type MediaAssetBinding struct {
	Id              uint64      `json:"id"              orm:"id"                description:""`                                                          //
	SceneCode       string      `json:"sceneCode"       orm:"scene_code"        description:""`                                                          //
	BizType         string      `json:"bizType"         orm:"biz_type"          description:""`                                                          //
	BizNo           string      `json:"bizNo"           orm:"biz_no"            description:""`                                                          //
	BindingField    string      `json:"bindingField"    orm:"binding_field"     description:"Slot name, e.g. main_images/detail_images/carousel_images"` // 插槽名，例如 main_images/detail_images/carousel_images。
	AssetId         uint64      `json:"assetId"         orm:"asset_id"          description:""`                                                          //
	SortOrder       int         `json:"sortOrder"       orm:"sort_order"        description:""`                                                          //
	IsActive        int         `json:"isActive"        orm:"is_active"         description:""`                                                          //
	ActiveAssetSlot uint64      `json:"activeAssetSlot" orm:"active_asset_slot" description:""`                                                          //
	ActiveSortSlot  int         `json:"activeSortSlot"  orm:"active_sort_slot"  description:""`                                                          //
	OperatorUserId  uint64      `json:"operatorUserId"  orm:"operator_user_id"  description:""`                                                          //
	RequestId       string      `json:"requestId"       orm:"request_id"        description:""`                                                          //
	UnboundAt       *gtime.Time `json:"unboundAt"       orm:"unbound_at"        description:""`                                                          //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        description:""`                                                          //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"        description:""`                                                          //
}
