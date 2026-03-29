package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
)

// Controller 閺?UserProfileService 閻?RPC 閹貉冨煑閸ｃ劌鐤勯悳鑸偓?
type Controller struct {
	v1.UnimplementedUserProfileServiceServer
}

// Register 鐏?controller 濞夈劌鍞介崚?gRPC 閺堝秴濮熼崳銊ｂ偓?
func Register(s *grpcx.GrpcServer) {
	v1.RegisterUserProfileServiceServer(s.Server, &Controller{})
}

// GetMyProfile 鏉烆剙褰傜拠閿嬬湴閸掔増婀囬崝鈥崇湴閼惧嘲褰囪ぐ鎾冲閻劍鍩涢悽璇插剼閵?
func (*Controller) GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error) {
	return service.UserProfile().GetMyProfile(ctx, req)
}

// UpdateMyProfile 鏉烆剙褰傜拠閿嬬湴閸掔増婀囬崝鈥崇湴閹笛嗩攽閻劍鍩涢悽璇插剼鐞涖儰绔甸妴?
func (*Controller) UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error) {
	return service.UserProfile().UpdateMyProfile(ctx, req)
}

// ListMyAddresses 鏉烆剙褰傜拠閿嬬湴閼惧嘲褰囪ぐ鎾冲閻劍鍩涢惃鍕勾閸р偓閸掓銆冮妴?
func (*Controller) ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error) {
	return service.UserProfile().ListMyAddresses(ctx, req)
}

// CreateMyAddress 鏉烆剙褰傜拠閿嬬湴閸掓稑缂撴稉鈧弶鈥虫勾閸р偓鐠佹澘缍嶉妴?
func (*Controller) CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error) {
	return service.UserProfile().CreateMyAddress(ctx, req)
}

// UpdateMyAddress 鏉烆剙褰傜拠閿嬬湴鐞涖儰绔甸幍鈧崷銊ユ勾閸р偓閵?
func (*Controller) UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error) {
	return service.UserProfile().UpdateMyAddress(ctx, req)
}

// DeleteMyAddress 鏉烆剙褰傜拠閿嬬湴閹笛嗩攽閸︽澘娼冩潪顖氬灩閵?
func (*Controller) DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error) {
	return service.UserProfile().DeleteMyAddress(ctx, req)
}

// SetMyDefaultAddress 鏉烆剙褰傜拠閿嬬湴鐠佸墽鐤嗛幋鏍ㄧ闂勩倝绮拋銈呮勾閸р偓閵?
func (*Controller) SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error) {
	return service.UserProfile().SetMyDefaultAddress(ctx, req)
}

// BatchGetProfileSummary 鏉烆剙褰傜拠閿嬬湴閹靛綊鍣虹拠璇插絿閻㈣鍎氶幗妯款洣閵?
func (*Controller) BatchGetProfileSummary(ctx context.Context, req *v1.BatchGetProfileSummaryReq) (res *v1.BatchGetProfileSummaryRes, err error) {
	return service.UserProfile().BatchGetProfileSummary(ctx, req)
}

// GetProfileByUserId 鏉烆剙褰傜拠閿嬬湴閹?user_id 閺屻儴顕楅悽璇插剼閵?
func (*Controller) GetProfileByUserId(ctx context.Context, req *v1.GetProfileByUserIdReq) (res *v1.GetProfileByUserIdRes, err error) {
	return service.UserProfile().GetProfileByUserId(ctx, req)
}

// ListAddressesByUserId 鏉烆剙褰傜拠閿嬬湴閺屻儴顕楅幐鍥х暰閻劍鍩涢惃鍕勾閸р偓閸掓銆冮妴?
func (*Controller) ListAddressesByUserId(ctx context.Context, req *v1.ListAddressesByUserIdReq) (res *v1.ListAddressesByUserIdRes, err error) {
	return service.UserProfile().ListAddressesByUserId(ctx, req)
}

// GetAddressById 鏉烆剙褰傜拠閿嬬湴閹?address_id 閺屻儴顕楅崷鏉挎絻閵?
func (*Controller) GetAddressById(ctx context.Context, req *v1.GetAddressByIdReq) (res *v1.GetAddressByIdRes, err error) {
	return service.UserProfile().GetAddressById(ctx, req)
}

// ReplaceMyAddress 鏉烆剙褰傜拠閿嬬湴閺囨寧宕查崷鏉挎絻楠炲墎鏁撻幋鎰煀鐞涘被鈧?
func (*Controller) ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error) {
	return service.UserProfile().ReplaceMyAddress(ctx, req)
}

// GetAddressSnapshotById 鏉烆剙褰傜拠閿嬬湴鐠囪褰囬崷鏉挎絻韫囶偆鍙庨妴?
func (*Controller) GetAddressSnapshotById(ctx context.Context, req *v1.GetAddressSnapshotByIdReq) (res *v1.GetAddressSnapshotByIdRes, err error) {
	return service.UserProfile().GetAddressSnapshotById(ctx, req)
}

// ResolveUserIdByPhone 鏉烆剙褰傜拠閿嬬湴闁俺绻冮幍瀣簚閸欏嘲寮介弻銉ф暏閹存灚鈧?
func (*Controller) ResolveUserIdByPhone(ctx context.Context, req *v1.ResolveUserIdByPhoneReq) (res *v1.ResolveUserIdByPhoneRes, err error) {
	return service.UserProfile().ResolveUserIdByPhone(ctx, req)
}

