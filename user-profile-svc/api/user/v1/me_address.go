package v1

import (
	pb "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

// ListMyAddressesReq defines HTTP request to list my addresses.
type ListMyAddressesReq struct {
	// Meta declares route + OpenAPI metadata for this endpoint.
	g.Meta `path:"/addresses" method:"get" tags:"User" summary:"List my addresses"`
	// IncludeDeleted controls whether deleted addresses are included.
	IncludeDeleted bool `json:"include_deleted" in:"query"`
}

// ListMyAddressesRes defines HTTP response for address list.
type ListMyAddressesRes struct {
	// Addresses is the address list.
	Addresses []*pb.UserAddress `json:"addresses"`
	// Page is current page.
	Page uint32 `json:"page"`
	// PageSize is current page size.
	PageSize uint32 `json:"page_size"`
	// Total is total count.
	Total uint32 `json:"total"`
	// AddressBookVersion is the current aggregate version for address-book updates.
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// CreateMyAddressReq defines HTTP request to create one address.
type CreateMyAddressReq struct {
	// Meta declares route + OpenAPI metadata for this endpoint.
	g.Meta `path:"/addresses" method:"post" tags:"User" summary:"Create my address"`
	// Address is create payload.
	Address *pb.AddressCreate `json:"address" v:"required#address is required"`
	// SetAsDefault indicates this address should become default.
	SetAsDefault bool `json:"set_as_default"`
	// ExpectedAddressBookVersion is suggested CAS version when set_as_default=true.
	ExpectedAddressBookVersion uint64 `json:"expected_address_book_version,omitempty"`
}

// CreateMyAddressRes defines HTTP response for created address.
type CreateMyAddressRes struct {
	// Address is created address.
	Address *pb.UserAddress `json:"address"`
	// AddressBookVersion is the updated aggregate version for address-book updates.
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// UpdateMyAddressReq defines HTTP request to patch an address.
type UpdateMyAddressReq struct {
	// Meta declares route + OpenAPI metadata for this endpoint.
	g.Meta `path:"/addresses/{addressId}" method:"patch" tags:"User" summary:"Patch my address"`
	// AddressId is the target address id.
	AddressId uint64 `json:"address_id" in:"path" v:"required#address_id is required"`
	// Address is patch payload.
	Address *pb.AddressPatch `json:"address" v:"required#address is required"`
	// UpdateMask tells which fields in Address should be applied.
	UpdateMask []string `json:"update_mask" v:"required#update_mask is required"`
	// ExpectedAddressVersion is CAS version read from latest address query.
	ExpectedAddressVersion uint64 `json:"expected_address_version" v:"required#expected_address_version is required"`
}

// UpdateMyAddressRes defines HTTP response for patched address.
type UpdateMyAddressRes struct {
	// Address is updated address snapshot.
	Address *pb.UserAddress `json:"address"`
	// AddressBookVersion is the updated aggregate version for address-book updates.
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// ReplaceMyAddressReq defines HTTP request to replace one address.
type ReplaceMyAddressReq struct {
	// Meta declares route + OpenAPI metadata for this endpoint.
	g.Meta `path:"/addresses/{sourceAddressId}/replace" method:"post" tags:"User" summary:"Replace my address"`
	// SourceAddressId is the source address id to be replaced.
	SourceAddressId uint64 `json:"source_address_id" in:"path" v:"required#source_address_id is required"`
	// Address is patch payload for the new address row.
	Address *pb.AddressPatch `json:"address" v:"required#address is required"`
	// UpdateMask tells which fields in Address should be applied.
	UpdateMask []string `json:"update_mask" v:"required#update_mask is required"`
	// SetAsDefault indicates new address should become default.
	SetAsDefault bool `json:"set_as_default,omitempty"`
	// ExpectedSourceAddressVersion is CAS version for source address.
	ExpectedSourceAddressVersion uint64 `json:"expected_source_address_version" v:"required#expected_source_address_version is required"`
	// ExpectedAddressBookVersion is suggested CAS version when set_as_default=true.
	ExpectedAddressBookVersion uint64 `json:"expected_address_book_version,omitempty"`
}

// ReplaceMyAddressRes defines HTTP response for replace operation.
type ReplaceMyAddressRes struct {
	// SourceAddressId is the retired source id.
	SourceAddressId uint64 `json:"source_address_id"`
	// NewAddress is the created replacement address.
	NewAddress *pb.UserAddress `json:"new_address"`
	// AddressBookVersion is the updated aggregate version for address-book updates.
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// DeleteMyAddressReq defines HTTP request to delete an address.
type DeleteMyAddressReq struct {
	// Meta declares route + OpenAPI metadata for this endpoint.
	g.Meta `path:"/addresses/{addressId}" method:"delete" tags:"User" summary:"Delete my address"`
	// AddressId is the target address id.
	AddressId uint64 `json:"address_id" in:"path" v:"required#address_id is required"`
	// ExpectedAddressVersion is CAS version read from latest address query.
	ExpectedAddressVersion uint64 `json:"expected_address_version" v:"required#expected_address_version is required"`
}

// DeleteMyAddressRes defines HTTP response for delete operation.
type DeleteMyAddressRes struct {
	// AddressBookVersion is the updated aggregate version for address-book updates.
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// SetMyDefaultAddressReq defines HTTP request to set an address as default.
type SetMyDefaultAddressReq struct {
	// Meta declares route + OpenAPI metadata for this endpoint.
	g.Meta `path:"/addresses/default" method:"post" tags:"User" summary:"Set or clear my default address"`
	// AddressId is optional; nil/0 means clear default address.
	AddressId *uint64 `json:"address_id,omitempty"`
	// ExpectedAddressBookVersion is CAS version read from latest list/get response.
	ExpectedAddressBookVersion uint64 `json:"expected_address_book_version" v:"required#expected_address_book_version is required"`
}

// SetMyDefaultAddressRes defines HTTP response for set default operation.
type SetMyDefaultAddressRes struct {
	// DefaultAddressId is current default address id, 0 means no default.
	DefaultAddressId uint64 `json:"default_address_id"`
	// AddressBookVersion is the updated aggregate version for address-book updates.
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}
