// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaAssetBinding is the golang structure of table media_asset_binding for DAO operations like Where/Data.
type MediaAssetBinding struct {
	g.Meta          `orm:"table:media_asset_binding, do:true"`
	Id              any         //
	SceneCode       any         //
	BizType         any         //
	BizNo           any         //
	BindingField    any         // Slot name, e.g. main_images/detail_images/carousel_images
	AssetId         any         //
	SortOrder       any         //
	IsActive        any         //
	ActiveAssetSlot any         //
	ActiveSortSlot  any         //
	OperatorUserId  any         //
	RequestId       any         //
	UnboundAt       *gtime.Time //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
