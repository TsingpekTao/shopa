// =================================================================================
// This file defines the HTTP controller interface for user APIs (v1).
// We keep it similar to the GoFrame-generated style so controllers stay consistent.
// =================================================================================

package user

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/user/v1"
)

// IUserV1 is the HTTP API surface for current-user profile/address management.
// Note: AuthN/AuthZ is handled by middleware/iam-svc, this service only uses user_id from ctx.
type IUserV1 interface {
	// GetMyProfile returns my profile, with optional addresses.
	GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error)
	// UpdateMyProfile patches my profile by update mask.
	UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error)

	// ListMyAddresses returns all addresses for current user.
	ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error)
	// CreateMyAddress creates an address.
	CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error)
	// ReplaceMyAddress replaces an address by creating a new row.
	ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error)
	// UpdateMyAddress patches an address.
	UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error)
	// DeleteMyAddress deletes an address (soft delete).
	DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error)
	// SetMyDefaultAddress sets or clears default address.
	SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error)
}
