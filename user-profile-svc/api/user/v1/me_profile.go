package v1

import (
	pb "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

// GetMyProfileReq defines HTTP request to get current user profile.
type GetMyProfileReq struct {
	// Meta declares route + OpenAPI metadata for this endpoint.
	g.Meta `path:"/profile" method:"get" tags:"User" summary:"Get my profile"`
	// IncludeAddresses controls whether to include address list in the same response.
	IncludeAddresses bool `json:"include_addresses" in:"query"`
}

// GetMyProfileRes defines HTTP response for my profile.
type GetMyProfileRes struct {
	// Profile is current user profile.
	Profile *pb.UserProfile `json:"profile"`
	// Addresses is optional address list.
	Addresses []*pb.UserAddress `json:"addresses,omitempty"`
	// AddressBookVersion is the current aggregate version for address-book updates.
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// UpdateMyProfileReq defines HTTP request to patch current user profile.
type UpdateMyProfileReq struct {
	// Meta declares route + OpenAPI metadata for this endpoint.
	g.Meta `path:"/profile" method:"patch" tags:"User" summary:"Patch my profile"`
	// Profile is the patch payload (values for fields selected by update_mask).
	Profile *pb.UserProfilePatch `json:"profile" v:"required#profile is required"`
	// UpdateMask tells which fields in Profile should be applied.
	UpdateMask []string `json:"update_mask" v:"required#update_mask is required"`
	// ExpectedProfileVersion is CAS version read from latest profile query.
	ExpectedProfileVersion uint64 `json:"expected_profile_version" v:"required#expected_profile_version is required"`
}

// UpdateMyProfileRes defines HTTP response for patched profile.
type UpdateMyProfileRes struct {
	// Profile is updated profile snapshot.
	Profile *pb.UserProfile `json:"profile"`
}
