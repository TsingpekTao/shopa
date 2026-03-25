// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MediaAssetBinding 是 media_asset_binding 表的 Go 结构体，供 DAO 的 Where/Data 等操作使用。
type MediaAssetBinding struct {
	g.Meta          `orm:"table:media_asset_binding, do:true"`
	Id              any         //
	SceneCode       any         //
	BizType         any         //
	BizNo           any         //
	BindingField    any         // 插槽名，例如 main_images/detail_images/carousel_images。
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
