package user

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	pb "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

// GetMyProfile 杩斿洖褰撳墠鐢ㄦ埛鐨勭敾鍍忓強鍙€夊湴鍧€銆?
func (*ControllerV1) GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error) {
	// 鐩存帴璋冪敤 gRPC 褰㈠紡鐨勬湇鍔″疄鐜帮紝閬垮厤閲嶅涓氬姟閫昏緫銆?
	out, err := service.UserProfile().GetMyProfile(ctx, &pb.GetMyProfileReq{
		IncludeAddresses: req.IncludeAddresses,
	})
	if err != nil {
		return nil, err
	}
	return &v1.GetMyProfileRes{
		Profile:            out.GetProfile(),
		Addresses:          out.GetAddresses(),
		AddressBookVersion: out.GetAddressBookVersion(),
	}, nil
}

// UpdateMyProfile 鏇存柊褰撳墠鐢ㄦ埛鐢诲儚骞惰繑鍥炴渶鏂扮粨鏋溿€?
func (*ControllerV1) UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error) {
	// 鏍￠獙璇锋眰杞借嵎涓嶄负 nil銆?
	if req.Profile == nil {
		return nil, gerror.New("profile is required")
	}
	// 杞彂鍒版湇鍔″眰浠ヤ繚鎸佷笟鍔￠€昏緫涓€鑷淬€?
	out, err := service.UserProfile().UpdateMyProfile(ctx, &pb.UpdateMyProfileReq{
		Profile:                req.Profile,
		UpdateMask:             req.UpdateMask,
		ExpectedProfileVersion: req.ExpectedProfileVersion,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateMyProfileRes{
		Profile: out.GetProfile(),
	}, nil
}

// ListMyAddresses 杩斿洖褰撳墠鐢ㄦ埛鐨勫湴鍧€鍒楄〃銆?
func (*ControllerV1) ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error) {
	// 鐩存帴璋冪敤鏈嶅姟灞傦紝褰撳墠 proto 宸茬Щ闄ゅ垎椤佃姹傚瓧娈点€?
	out, err := service.UserProfile().ListMyAddresses(ctx, &pb.ListMyAddressesReq{
		IncludeDeleted: req.IncludeDeleted,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ListMyAddressesRes{
		Addresses:          out.GetAddresses(),
		Page:               out.GetPage(),
		PageSize:           out.GetPageSize(),
		Total:              out.GetTotal(),
		AddressBookVersion: out.GetAddressBookVersion(),
	}, nil
}

// CreateMyAddress 鍒涘缓涓€鏉＄敤鎴峰湴鍧€銆?
func (*ControllerV1) CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error) {
	// 鏍￠獙鍦板潃杞借嵎涓嶄负 nil銆?
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// 杞彂鍒版湇鍔″眰浠ュ鐢ㄤ簨鍔″拰榛樿鍦板潃閫昏緫銆?
	out, err := service.UserProfile().CreateMyAddress(ctx, &pb.CreateMyAddressReq{
		Address:                    req.Address,
		SetAsDefault:               req.SetAsDefault,
		ExpectedAddressBookVersion: req.ExpectedAddressBookVersion,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateMyAddressRes{
		Address:            out.GetAddress(),
		AddressBookVersion: out.GetAddressBookVersion(),
	}, nil
}

// UpdateMyAddress 鏇存柊鎸囧畾鐨勭敤鎴峰湴鍧€銆?
func (*ControllerV1) UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error) {
	// 鏍￠獙鍦板潃杞借嵎涓嶄负 nil銆?
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// 杞彂鍒版湇鍔″眰浠ヤ繚鎸佷竴鑷寸殑瑙勫垯澶勭悊銆?
	out, err := service.UserProfile().UpdateMyAddress(ctx, &pb.UpdateMyAddressReq{
		AddressId:              req.AddressId,
		Address:                req.Address,
		UpdateMask:             req.UpdateMask,
		ExpectedAddressVersion: req.ExpectedAddressVersion,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateMyAddressRes{
		Address:            out.GetAddress(),
		AddressBookVersion: out.GetAddressBookVersion(),
	}, nil
}

// ReplaceMyAddress 鏂板缓涓€鏉″湴鍧€鏇挎崲鍘熷湴鍧€骞舵洿鏂伴粯璁ゆ€佸娍銆?
func (*ControllerV1) ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error) {
	// 鏍￠獙鍦板潃杞借嵎涓嶄负 nil銆?
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// 杞彂鍒版湇鍔″眰浠ュ鐢ㄦ浛鎹簨鍔￠€昏緫銆?
	out, err := service.UserProfile().ReplaceMyAddress(ctx, &pb.ReplaceMyAddressReq{
		SourceAddressId:              req.SourceAddressId,
		Address:                      req.Address,
		UpdateMask:                   req.UpdateMask,
		SetAsDefault:                 req.SetAsDefault,
		ExpectedSourceAddressVersion: req.ExpectedSourceAddressVersion,
		ExpectedAddressBookVersion:   req.ExpectedAddressBookVersion,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ReplaceMyAddressRes{
		SourceAddressId:    out.GetSourceAddressId(),
		NewAddress:         out.GetNewAddress(),
		AddressBookVersion: out.GetAddressBookVersion(),
	}, nil
}

// DeleteMyAddress 鍒犻櫎鎸囧畾鍦板潃骞朵繚鐣欑増鏈拷韪€?
func (*ControllerV1) DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error) {
	// 杞彂鍒版湇鍔″眰浠ュ鐢ㄥ綊灞炴牎楠屼笌杞垹琛屼负銆?
	out, err := service.UserProfile().DeleteMyAddress(ctx, &pb.DeleteMyAddressReq{
		AddressId:              req.AddressId,
		ExpectedAddressVersion: req.ExpectedAddressVersion,
	})
	if err != nil {
		return nil, err
	}
	return &v1.DeleteMyAddressRes{
		AddressBookVersion: out.GetAddressBookVersion(),
	}, nil
}

// SetMyDefaultAddress 璁剧疆鎴栨竻闄ょ敤鎴风殑榛樿鍦板潃銆?
func (*ControllerV1) SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error) {
	// 杞彂鍒版湇鍔″眰浠ヤ繚璇侀粯璁ゅ湴鍧€鍞竴銆?
	out, err := service.UserProfile().SetMyDefaultAddress(ctx, &pb.SetMyDefaultAddressReq{
		AddressId:                  req.AddressId,
		ExpectedAddressBookVersion: req.ExpectedAddressBookVersion,
	})
	if err != nil {
		return nil, err
	}
	return &v1.SetMyDefaultAddressRes{
		DefaultAddressId:   out.GetDefaultAddressId(),
		AddressBookVersion: out.GetAddressBookVersion(),
	}, nil
}

