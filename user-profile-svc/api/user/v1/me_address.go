package v1

import (
	pb "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

// ListMyAddressesReq 定义列出当前用户地址的 HTTP 请求。
type ListMyAddressesReq struct {
	// Meta 声明此端点的路由与 OpenAPI 元数据。
	g.Meta `path:"/addresses" method:"get" tags:"User" summary:"List my addresses"`
	// IncludeDeleted 控制是否包含已删除的地址。
	IncludeDeleted bool `json:"include_deleted" in:"query"`
}

// ListMyAddressesRes 定义地址列表的 HTTP 响应。
type ListMyAddressesRes struct {
	// Addresses 是地址列表。
	Addresses []*pb.UserAddress `json:"addresses"`
	// Page 表示当前页号。
	Page uint32 `json:"page"`
	// PageSize 表示当前页大小。
	PageSize uint32 `json:"page_size"`
	// Total 表示总记录数。
	Total uint32 `json:"total"`
	// AddressBookVersion 表示地址簿的汇总版本。
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// CreateMyAddressReq 定义创建地址的 HTTP 请求。
type CreateMyAddressReq struct {
	// Meta 声明此端点的路由与 OpenAPI 元数据。
	g.Meta `path:"/addresses" method:"post" tags:"User" summary:"Create my address"`
	// Address 是创建地址的入参。
	Address *pb.AddressCreate `json:"address" v:"required#address is required"`
	// SetAsDefault 表示此地址应设为默认。
	SetAsDefault bool `json:"set_as_default"`
	// ExpectedAddressBookVersion 表示在 set_as_default=true 时建议的地址簿 CAS 版本。
	ExpectedAddressBookVersion uint64 `json:"expected_address_book_version,omitempty"`
}

// CreateMyAddressRes 定义创建地址后的 HTTP 响应。
type CreateMyAddressRes struct {
	// Address 是创建后的地址。
	Address *pb.UserAddress `json:"address"`
	// AddressBookVersion 表示最新地址簿版本。
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// UpdateMyAddressReq 定义更新地址的 HTTP 请求。
type UpdateMyAddressReq struct {
	// Meta 声明此端点的路由与 OpenAPI 元数据。
	g.Meta `path:"/addresses/{addressId}" method:"patch" tags:"User" summary:"Patch my address"`
	// AddressId 是目标地址 ID。
	AddressId uint64 `json:"address_id" in:"path" v:"required#address_id is required"`
	// Address 是用于补丁的地址数据。
	Address *pb.AddressPatch `json:"address" v:"required#address is required"`
	// UpdateMask 指示需要更新的字段。
	UpdateMask []string `json:"update_mask" v:"required#update_mask is required"`
	// ExpectedAddressVersion 是从最新地址查询读取的 CAS 版本。
	ExpectedAddressVersion uint64 `json:"expected_address_version" v:"required#expected_address_version is required"`
}

// UpdateMyAddressRes 定义地址更新后的 HTTP 响应。
type UpdateMyAddressRes struct {
	// Address 是更新后的地址快照。
	Address *pb.UserAddress `json:"address"`
	// AddressBookVersion 表示最新地址簿版本。
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// ReplaceMyAddressReq 定义替换地址的 HTTP 请求。
type ReplaceMyAddressReq struct {
	// Meta 声明此端点的路由与 OpenAPI 元数据。
	g.Meta `path:"/addresses/{sourceAddressId}/replace" method:"post" tags:"User" summary:"Replace my address"`
	// SourceAddressId 是要替换的源地址 ID。
	SourceAddressId uint64 `json:"source_address_id" in:"path" v:"required#source_address_id is required"`
	// Address 是用于新地址行的补丁数据。
	Address *pb.AddressPatch `json:"address" v:"required#address is required"`
	// UpdateMask 指示需应用的新地址字段。
	UpdateMask []string `json:"update_mask" v:"required#update_mask is required"`
	// SetAsDefault 表示新地址应设为默认。
	SetAsDefault bool `json:"set_as_default,omitempty"`
	// ExpectedSourceAddressVersion 是源地址的 CAS 版本。
	ExpectedSourceAddressVersion uint64 `json:"expected_source_address_version" v:"required#expected_source_address_version is required"`
	// ExpectedAddressBookVersion 表示当 set_as_default=true 时的建议地址簿版本。
	ExpectedAddressBookVersion uint64 `json:"expected_address_book_version,omitempty"`
}

// ReplaceMyAddressRes 定义替换操作的 HTTP 响应。
type ReplaceMyAddressRes struct {
	// SourceAddressId 是已退役的源地址 ID。
	SourceAddressId uint64 `json:"source_address_id"`
	// NewAddress 是新建的替代地址。
	NewAddress *pb.UserAddress `json:"new_address"`
	// AddressBookVersion 表示最新地址簿版本。
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// DeleteMyAddressReq 定义删除地址的 HTTP 请求。
type DeleteMyAddressReq struct {
	// Meta 声明此端点的路由与 OpenAPI 元数据。
	g.Meta `path:"/addresses/{addressId}" method:"delete" tags:"User" summary:"Delete my address"`
	// AddressId 是目标地址 ID。
	AddressId uint64 `json:"address_id" in:"path" v:"required#address_id is required"`
	// ExpectedAddressVersion 是从最新地址查询读取的 CAS 版本。
	ExpectedAddressVersion uint64 `json:"expected_address_version" v:"required#expected_address_version is required"`
}

// DeleteMyAddressRes 定义删除操作的 HTTP 响应。
type DeleteMyAddressRes struct {
	// AddressBookVersion 表示最新地址簿版本。
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}

// SetMyDefaultAddressReq 定义设置默认地址的 HTTP 请求。
type SetMyDefaultAddressReq struct {
	// Meta 声明此端点的路由与 OpenAPI 元数据。
	g.Meta `path:"/addresses/default" method:"post" tags:"User" summary:"Set or clear my default address"`
	// AddressId 可选；nil/0 表示清除默认地址。
	AddressId *uint64 `json:"address_id,omitempty"`
	// ExpectedAddressBookVersion 是从最近列表/查询响应读取的地址簿 CAS 版本。
	ExpectedAddressBookVersion uint64 `json:"expected_address_book_version" v:"required#expected_address_book_version is required"`
}

// SetMyDefaultAddressRes 定义设置默认地址操作的 HTTP 响应。
type SetMyDefaultAddressRes struct {
	// DefaultAddressId 是当前默认地址 ID，0 表示没有默认地址。
	DefaultAddressId uint64 `json:"default_address_id"`
	// AddressBookVersion 表示最新地址簿版本。
	AddressBookVersion uint64 `json:"address_book_version,omitempty"`
}
