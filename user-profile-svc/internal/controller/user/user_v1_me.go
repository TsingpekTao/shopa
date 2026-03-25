package user

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/user/v1"
	pb "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

// GetMyProfile 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (*ControllerV1) GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error) {
	// Delegate to gRPC-shaped service method to avoid duplicating business rules.
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

// UpdateMyProfile 鎸夋潯浠舵洿鏂版暟鎹苟杩斿洖鏈€鏂扮粨鏋溿€
func (*ControllerV1) UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error) {
	// Protect against nil pointer payload.
	if req.Profile == nil {
		return nil, gerror.New("profile is required")
	}
	// Convert update mask array into protobuf FieldMask.
	mask := &fieldmaskpb.FieldMask{
		Paths: req.UpdateMask,
	}
	// Delegate to service to keep logic consistent with RPC behavior.
	out, err := service.UserProfile().UpdateMyProfile(ctx, &pb.UpdateMyProfileReq{
		Profile:                req.Profile,
		UpdateMask:             mask,
		ExpectedProfileVersion: req.ExpectedProfileVersion,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateMyProfileRes{
		Profile: out.GetProfile(),
	}, nil
}

// ListMyAddresses 鎸夋潯浠惰鍙栧苟杩斿洖鍒楄〃缁撴灉銆
func (*ControllerV1) ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error) {
	// Delegate to service. New proto removed paging in request.
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

// CreateMyAddress 鍒涘缓鏂拌褰曞苟杩斿洖鍒涘缓缁撴灉銆
func (*ControllerV1) CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error) {
	// Protect against nil pointer payload.
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// Delegate to service for transaction + default logic.
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

// UpdateMyAddress 鎸夋潯浠舵洿鏂版暟鎹苟杩斿洖鏈€鏂扮粨鏋溿€
func (*ControllerV1) UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error) {
	// Protect against nil pointer payload.
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// Convert update mask array into protobuf FieldMask.
	mask := &fieldmaskpb.FieldMask{
		Paths: req.UpdateMask,
	}
	// Delegate to service to keep rules consistent.
	out, err := service.UserProfile().UpdateMyAddress(ctx, &pb.UpdateMyAddressReq{
		AddressId:              req.AddressId,
		Address:                req.Address,
		UpdateMask:             mask,
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

// ReplaceMyAddress 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (*ControllerV1) ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error) {
	// Protect against nil pointer payload.
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// Convert update mask array into protobuf FieldMask.
	mask := &fieldmaskpb.FieldMask{
		Paths: req.UpdateMask,
	}
	// Delegate to service for replacement transaction.
	out, err := service.UserProfile().ReplaceMyAddress(ctx, &pb.ReplaceMyAddressReq{
		SourceAddressId:              req.SourceAddressId,
		Address:                      req.Address,
		UpdateMask:                   mask,
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

// DeleteMyAddress 鎵ц鍒犻櫎娴佺▼骞惰繑鍥炲鐞嗙粨鏋溿€
func (*ControllerV1) DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error) {
	// Delegate to service for ownership checks and soft delete behavior.
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

// SetMyDefaultAddress 璁剧疆鐘舵€佹垨榛樿鍊煎苟淇濊瘉绾︽潫涓€鑷淬€
func (*ControllerV1) SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error) {
	// Delegate to service for transaction ensuring exactly one default.
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
