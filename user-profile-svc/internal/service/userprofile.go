// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 说明：如手动维护此接口文件，可删除这些提示注释。
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
)

type (
	IUserProfile interface {
		// ApplyRegisterInitEvent 消费注册初始化事件（事件幂等 + 时序保护 + 昵称来源保护）。
		ApplyRegisterInitEvent(ctx context.Context, event RegisterInitEvent) error
		// GetMyProfile 查询“当前登录用户”资料，可按需附带地址列表。
		GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error)
		// UpdateMyProfile 按 FieldMask 局部更新当前用户资料，并推进 profile_version。
		UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error)
		// ListMyAddresses 查询当前用户地址簿，可选择包含已删除地址。
		ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error)
		// CreateMyAddress 创建新地址，并在事务内处理默认地址投影。
		CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error)
		// UpdateMyAddress 按 FieldMask 原地更新一条 ACTIVE 地址。
		UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error)
		// DeleteMyAddress 软删除地址，并在必要时清理 profile.default_address_id。
		DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error)
		// SetMyDefaultAddress 设置或清空默认地址，并推进 address_book_version。
		SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error)
		// BatchGetProfileSummary 批量返回轻量资料摘要（供内部服务调用）。
		BatchGetProfileSummary(ctx context.Context, req *v1.BatchGetProfileSummaryReq) (res *v1.BatchGetProfileSummaryRes, err error)
		// GetProfileByUserId 按显式 user_id 查询资料（内部接口）。
		GetProfileByUserId(ctx context.Context, req *v1.GetProfileByUserIdReq) (res *v1.GetProfileByUserIdRes, err error)
		// ListAddressesByUserId 按显式 user_id 查询地址列表（内部接口）。
		ListAddressesByUserId(ctx context.Context, req *v1.ListAddressesByUserIdReq) (res *v1.ListAddressesByUserIdRes, err error)
		// GetAddressById 按 address_id 查询单条地址（内部接口）。
		GetAddressById(ctx context.Context, req *v1.GetAddressByIdReq) (res *v1.GetAddressByIdRes, err error)
		// ReplaceMyAddress 使用“新增新行 + 旧行置 REPLACED”语义替换地址，避免下游快照污染。
		ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error)
		// GetAddressSnapshotById 读取可用于订单冗余的地址快照（携带源地址 ID/版本）。
		GetAddressSnapshotById(ctx context.Context, req *v1.GetAddressSnapshotByIdReq) (res *v1.GetAddressSnapshotByIdRes, err error)
		// ResolveUserIdByPhone 按手机号反查用户（当前实现基于 ACTIVE 地址表 receiver_phone），并返回脱敏手机号。
		ResolveUserIdByPhone(ctx context.Context, req *v1.ResolveUserIdByPhoneReq) (res *v1.ResolveUserIdByPhoneRes, err error)
	}
)

var (
	localUserProfile IUserProfile
)

// UserProfile 返回 IUserProfile 的当前注册实现。
// 若未在逻辑层执行 RegisterUserProfile，会直接 panic 提醒启动配置缺失。
func UserProfile() IUserProfile {
	if localUserProfile == nil {
		panic("implement not found for interface IUserProfile, forgot register?")
	}
	return localUserProfile
}

// RegisterUserProfile 注入 IUserProfile 的具体实现（通常在 logic 包 init 中调用）。
func RegisterUserProfile(i IUserProfile) {
	localUserProfile = i
}
