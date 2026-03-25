package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
)

// Controller 是 UserProfileService 的 RPC 控制器实现。
type Controller struct {
	v1.UnimplementedUserProfileServiceServer
}

// Register 将 controller 注册到 gRPC 服务器。
func Register(s *grpcx.GrpcServer) {
	v1.RegisterUserProfileServiceServer(s.Server, &Controller{})
}

// GetMyProfile 转发请求到服务层获取当前用户画像。
func (*Controller) GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error) {
	return service.UserProfile().GetMyProfile(ctx, req)
}

// UpdateMyProfile 转发请求到服务层执行用户画像补丁。
func (*Controller) UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error) {
	return service.UserProfile().UpdateMyProfile(ctx, req)
}

// ListMyAddresses 转发请求获取当前用户的地址列表。
func (*Controller) ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error) {
	return service.UserProfile().ListMyAddresses(ctx, req)
}

// CreateMyAddress 转发请求创建一条地址记录。
func (*Controller) CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error) {
	return service.UserProfile().CreateMyAddress(ctx, req)
}

// UpdateMyAddress 转发请求补丁所在地址。
func (*Controller) UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error) {
	return service.UserProfile().UpdateMyAddress(ctx, req)
}

// DeleteMyAddress 转发请求执行地址软删。
func (*Controller) DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error) {
	return service.UserProfile().DeleteMyAddress(ctx, req)
}

// SetMyDefaultAddress 转发请求设置或清除默认地址。
func (*Controller) SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error) {
	return service.UserProfile().SetMyDefaultAddress(ctx, req)
}

// BatchGetProfileSummary 转发请求批量读取画像摘要。
func (*Controller) BatchGetProfileSummary(ctx context.Context, req *v1.BatchGetProfileSummaryReq) (res *v1.BatchGetProfileSummaryRes, err error) {
	return service.UserProfile().BatchGetProfileSummary(ctx, req)
}

// GetProfileByUserId 转发请求按 user_id 查询画像。
func (*Controller) GetProfileByUserId(ctx context.Context, req *v1.GetProfileByUserIdReq) (res *v1.GetProfileByUserIdRes, err error) {
	return service.UserProfile().GetProfileByUserId(ctx, req)
}

// ListAddressesByUserId 转发请求查询指定用户的地址列表。
func (*Controller) ListAddressesByUserId(ctx context.Context, req *v1.ListAddressesByUserIdReq) (res *v1.ListAddressesByUserIdRes, err error) {
	return service.UserProfile().ListAddressesByUserId(ctx, req)
}

// GetAddressById 转发请求按 address_id 查询地址。
func (*Controller) GetAddressById(ctx context.Context, req *v1.GetAddressByIdReq) (res *v1.GetAddressByIdRes, err error) {
	return service.UserProfile().GetAddressById(ctx, req)
}

// ReplaceMyAddress 转发请求替换地址并生成新行。
func (*Controller) ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error) {
	return service.UserProfile().ReplaceMyAddress(ctx, req)
}

// GetAddressSnapshotById 转发请求读取地址快照。
func (*Controller) GetAddressSnapshotById(ctx context.Context, req *v1.GetAddressSnapshotByIdReq) (res *v1.GetAddressSnapshotByIdRes, err error) {
	return service.UserProfile().GetAddressSnapshotById(ctx, req)
}

// ResolveUserIdByPhone 转发请求通过手机号反查用户。
func (*Controller) ResolveUserIdByPhone(ctx context.Context, req *v1.ResolveUserIdByPhoneReq) (res *v1.ResolveUserIdByPhoneRes, err error) {
	return service.UserProfile().ResolveUserIdByPhone(ctx, req)
}
