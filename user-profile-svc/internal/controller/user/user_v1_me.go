package user

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/user/v1"
	pb "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

// GetMyProfile 返回当前用户的画像及可选地址。
func (*ControllerV1) GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error) {
	// 直接调用 gRPC 形式的服务实现，避免重复业务逻辑。
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

// UpdateMyProfile 更新当前用户画像并返回最新结果。
func (*ControllerV1) UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error) {
	// 校验请求载荷不为 nil。
	if req.Profile == nil {
		return nil, gerror.New("profile is required")
	}
	// 将字段掩码转换为 protobuf FieldMask。
	mask := &fieldmaskpb.FieldMask{
		Paths: req.UpdateMask,
	}
	// 转发到服务层以保持业务逻辑一致。
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

// ListMyAddresses 返回当前用户的地址列表。
func (*ControllerV1) ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error) {
	// 直接调用服务层，当前 proto 已移除分页请求字段。
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

// CreateMyAddress 创建一条用户地址。
func (*ControllerV1) CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error) {
	// 校验地址载荷不为 nil。
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// 转发到服务层以复用事务和默认地址逻辑。
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

// UpdateMyAddress 更新指定的用户地址。
func (*ControllerV1) UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error) {
	// 校验地址载荷不为 nil。
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// 将字段掩码转换为 protobuf FieldMask。
	mask := &fieldmaskpb.FieldMask{
		Paths: req.UpdateMask,
	}
	// 转发到服务层以保持一致的规则处理。
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

// ReplaceMyAddress 新建一条地址替换原地址并更新默认态势。
func (*ControllerV1) ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error) {
	// 校验地址载荷不为 nil。
	if req.Address == nil {
		return nil, gerror.New("address is required")
	}
	// 将字段掩码转换为 protobuf FieldMask。
	mask := &fieldmaskpb.FieldMask{
		Paths: req.UpdateMask,
	}
	// 转发到服务层以复用替换事务逻辑。
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

// DeleteMyAddress 删除指定地址并保留版本追踪。
func (*ControllerV1) DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error) {
	// 转发到服务层以复用归属校验与软删行为。
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

// SetMyDefaultAddress 设置或清除用户的默认地址。
func (*ControllerV1) SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error) {
	// 转发到服务层以保证默认地址唯一。
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
