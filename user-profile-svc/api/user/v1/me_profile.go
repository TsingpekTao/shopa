package v1

import (
	pb "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

// GetMyProfileReq 定义获取当前用户画像的 HTTP 请求。
type GetMyProfileReq struct {
	// Meta 声明此端点的路由与 OpenAPI 元数据。
	g.Meta `path:"/profile" method:"get" tags:"User" summary:"Get my profile"`
	// IncludeAddresses 控制是否在同一次响应中返回地址列表。
	IncludeAddresses bool `json:"include_addresses" in:"query"`
}

// GetMyProfileRes 定义用户画像的 HTTP 响应。
type GetMyProfileRes struct {
	// Profile 是当前用户画像。
	Profile *pb.UserProfile `json:"profile"`
	// Addresses 是可选的地址列表。
	Addresses []*pb.UserAddress `json:"addresses,omitempty"`
	// AddressBookVersion 表示地址簿汇总版本。
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// UpdateMyProfileReq 定义更新当前用户画像的 HTTP 请求。
type UpdateMyProfileReq struct {
	// Meta 声明此端点的路由与 OpenAPI 元数据。
	g.Meta `path:"/profile" method:"patch" tags:"User" summary:"Patch my profile"`
	// Profile 是补丁载荷（包含 update_mask 选中字段的值）。
	Profile *pb.UserProfilePatch `json:"profile" v:"required#profile is required"`
	// UpdateMask 指示需要应用的字段。
	UpdateMask []string `json:"update_mask" v:"required#update_mask is required"`
	// ExpectedProfileVersion 是从最近画像查询读取的 CAS 版本。
	ExpectedProfileVersion uint64 `json:"expected_profile_version" v:"required#expected_profile_version is required"`
}

// UpdateMyProfileRes 定义更新画像后的 HTTP 响应。
type UpdateMyProfileRes struct {
	// Profile 是更新后的画像快照。
	Profile *pb.UserProfile `json:"profile"`
}
