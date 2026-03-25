// =================================================================================
// 本文件定义用户 API (v1) 的 HTTP 控制器接口。
// 保持与 GoFrame 生成风格一致，确保控制器结构统一。
// =================================================================================

package user

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/user/v1"
)

// IUserV1 是当前用户画像/地址管理的 HTTP API 接口。
// 注意：认证/授权由 middleware/iam-svc 负责，本服务仅利用上下文中的 user_id。
type IUserV1 interface {
	// GetMyProfile 返回当前用户画像，可选包含地址列表。
	GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error)
	// UpdateMyProfile 按更新掩码补丁当前用户画像。
	UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error)

	// ListMyAddresses 列出当前用户的所有地址。
	ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error)
	// CreateMyAddress 创建一条用户地址。
	CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error)
	// ReplaceMyAddress 通过新建行替换指定地址。
	ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error)
	// UpdateMyAddress 补丁指定地址。
	UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error)
	// DeleteMyAddress 软删除指定地址。
	DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error)
	// SetMyDefaultAddress 设置或清空默认地址。
	SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error)
}
