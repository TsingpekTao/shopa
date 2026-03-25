package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
)

// Controller is the RPC controller implementation for UserProfileService.
type Controller struct {
	v1.UnimplementedUserProfileServiceServer
}

// Register 澶勭悊娉ㄥ唽涓绘祦绋嬪強鍒濆鍖栧姩浣溿€
func Register(s *grpcx.GrpcServer) {
	v1.RegisterUserProfileServiceServer(s.Server, &Controller{})
}

// GetMyProfile 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (*Controller) GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error) {
	return service.UserProfile().GetMyProfile(ctx, req)
}

// UpdateMyProfile 鎸夋潯浠舵洿鏂版暟鎹苟杩斿洖鏈€鏂扮粨鏋溿€
func (*Controller) UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error) {
	return service.UserProfile().UpdateMyProfile(ctx, req)
}

// ListMyAddresses 鎸夋潯浠惰鍙栧苟杩斿洖鍒楄〃缁撴灉銆
func (*Controller) ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error) {
	return service.UserProfile().ListMyAddresses(ctx, req)
}

// CreateMyAddress 鍒涘缓鏂拌褰曞苟杩斿洖鍒涘缓缁撴灉銆
func (*Controller) CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error) {
	return service.UserProfile().CreateMyAddress(ctx, req)
}

// UpdateMyAddress 鎸夋潯浠舵洿鏂版暟鎹苟杩斿洖鏈€鏂扮粨鏋溿€
func (*Controller) UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error) {
	return service.UserProfile().UpdateMyAddress(ctx, req)
}

// DeleteMyAddress 鎵ц鍒犻櫎娴佺▼骞惰繑鍥炲鐞嗙粨鏋溿€
func (*Controller) DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error) {
	return service.UserProfile().DeleteMyAddress(ctx, req)
}

// SetMyDefaultAddress 璁剧疆鐘舵€佹垨榛樿鍊煎苟淇濊瘉绾︽潫涓€鑷淬€
func (*Controller) SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error) {
	return service.UserProfile().SetMyDefaultAddress(ctx, req)
}

// BatchGetProfileSummary 鎵归噺澶勭悊璇锋眰锛屽噺灏戝線杩斿紑閿€銆
func (*Controller) BatchGetProfileSummary(ctx context.Context, req *v1.BatchGetProfileSummaryReq) (res *v1.BatchGetProfileSummaryRes, err error) {
	return service.UserProfile().BatchGetProfileSummary(ctx, req)
}

// GetProfileByUserId 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (*Controller) GetProfileByUserId(ctx context.Context, req *v1.GetProfileByUserIdReq) (res *v1.GetProfileByUserIdRes, err error) {
	return service.UserProfile().GetProfileByUserId(ctx, req)
}

// ListAddressesByUserId 鎸夋潯浠惰鍙栧苟杩斿洖鍒楄〃缁撴灉銆
func (*Controller) ListAddressesByUserId(ctx context.Context, req *v1.ListAddressesByUserIdReq) (res *v1.ListAddressesByUserIdRes, err error) {
	return service.UserProfile().ListAddressesByUserId(ctx, req)
}

// GetAddressById 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (*Controller) GetAddressById(ctx context.Context, req *v1.GetAddressByIdReq) (res *v1.GetAddressByIdRes, err error) {
	return service.UserProfile().GetAddressById(ctx, req)
}

// ReplaceMyAddress 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (*Controller) ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error) {
	return service.UserProfile().ReplaceMyAddress(ctx, req)
}

// GetAddressSnapshotById 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (*Controller) GetAddressSnapshotById(ctx context.Context, req *v1.GetAddressSnapshotByIdReq) (res *v1.GetAddressSnapshotByIdRes, err error) {
	return service.UserProfile().GetAddressSnapshotById(ctx, req)
}

// ResolveUserIdByPhone 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (*Controller) ResolveUserIdByPhone(ctx context.Context, req *v1.ResolveUserIdByPhoneReq) (res *v1.ResolveUserIdByPhoneRes, err error) {
	return service.UserProfile().ResolveUserIdByPhone(ctx, req)
}
