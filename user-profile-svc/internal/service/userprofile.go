// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
)

type (
	IUserProfile interface {
		// ApplyRegisterInitEvent applies register-init event idempotently.
		ApplyRegisterInitEvent(ctx context.Context, event RegisterInitEvent) error
		// GetMyProfile returns current user's profile and optional address list.
		GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error)
		// UpdateMyProfile updates the current profile by field mask.
		UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error)
		// ListMyAddresses returns address list for current user.
		ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error)
		// CreateMyAddress creates a new address for current user.
		CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error)
		// UpdateMyAddress updates an existing address for current user.
		UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error)
		// DeleteMyAddress soft-deletes one address for current user.
		DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error)
		// SetMyDefaultAddress sets default address for current user.
		SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error)
		// BatchGetProfileSummary returns light-weight profile info for internal callers.
		BatchGetProfileSummary(ctx context.Context, req *v1.BatchGetProfileSummaryReq) (res *v1.BatchGetProfileSummaryRes, err error)
		// GetProfileByUserId returns profile for an explicit user id (internal use).
		GetProfileByUserId(ctx context.Context, req *v1.GetProfileByUserIdReq) (res *v1.GetProfileByUserIdRes, err error)
		// ListAddressesByUserId returns all addresses for an explicit user id (internal use).
		ListAddressesByUserId(ctx context.Context, req *v1.ListAddressesByUserIdReq) (res *v1.ListAddressesByUserIdRes, err error)
		// GetAddressById returns one address by address id (internal use).
		GetAddressById(ctx context.Context, req *v1.GetAddressByIdReq) (res *v1.GetAddressByIdRes, err error)
		// ReplaceMyAddress replaces an existing address by creating a new row and retiring the old row.
		ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error)
		// GetAddressSnapshotById returns order-safe immutable snapshot payload from an address row.
		GetAddressSnapshotById(ctx context.Context, req *v1.GetAddressSnapshotByIdReq) (res *v1.GetAddressSnapshotByIdRes, err error)
		// ResolveUserIdByPhone resolves user by phone for internal tools.
		// Current implementation searches receiver_phone in active addresses as a pragmatic fallback.
		ResolveUserIdByPhone(ctx context.Context, req *v1.ResolveUserIdByPhoneReq) (res *v1.ResolveUserIdByPhoneRes, err error)
	}
)

var (
	localUserProfile IUserProfile
)

// UserProfile 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func UserProfile() IUserProfile {
	if localUserProfile == nil {
		panic("implement not found for interface IUserProfile, forgot register?")
	}
	return localUserProfile
}

// RegisterUserProfile 澶勭悊娉ㄥ唽涓绘祦绋嬪強鍒濆鍖栧姩浣溿€
func RegisterUserProfile(i IUserProfile) {
	localUserProfile = i
}
