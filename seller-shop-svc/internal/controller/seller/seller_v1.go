package seller

import (
	"context"

	sellerv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/seller/v1"
	pb "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
)

func (c *ControllerV1) CreateApplicationDraft(ctx context.Context, req *sellerv1.CreateApplicationDraftReq) (res *sellerv1.CreateApplicationDraftRes, err error) {
	out, err := c.svc.CreateApplicationDraft(withRequestMetadata(ctx), &pb.CreateApplicationDraftReq{
		Entity: req.Entity,
		Shop:   req.Shop,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.CreateApplicationDraftRes{Application: out.GetApplication()}, nil
}

func (c *ControllerV1) UpdateApplicationDraft(ctx context.Context, req *sellerv1.UpdateApplicationDraftReq) (res *sellerv1.UpdateApplicationDraftRes, err error) {
	out, err := c.svc.UpdateApplicationDraft(withRequestMetadata(ctx), &pb.UpdateApplicationDraftReq{
		ApplicationNo:   req.ApplicationNo,
		ExpectedVersion: req.ExpectedVersion,
		Entity:          req.Entity,
		Shop:            req.Shop,
		UpdateMask:      toFieldMask(req.UpdateMask),
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.UpdateApplicationDraftRes{Application: out.GetApplication()}, nil
}

func (c *ControllerV1) SubmitApplication(ctx context.Context, req *sellerv1.SubmitApplicationReq) (res *sellerv1.SubmitApplicationRes, err error) {
	out, err := c.svc.SubmitApplication(withRequestMetadata(ctx), &pb.SubmitApplicationReq{
		ApplicationNo:   req.ApplicationNo,
		ExpectedVersion: req.ExpectedVersion,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.SubmitApplicationRes{Application: out.GetApplication()}, nil
}

func (c *ControllerV1) ResubmitApplication(ctx context.Context, req *sellerv1.ResubmitApplicationReq) (res *sellerv1.ResubmitApplicationRes, err error) {
	out, err := c.svc.ResubmitApplication(withRequestMetadata(ctx), &pb.ResubmitApplicationReq{
		RejectedApplicationNo: req.RejectedApplicationNo,
		ExpectedVersion:       req.ExpectedVersion,
		Entity:                req.Entity,
		Shop:                  req.Shop,
		UpdateMask:            toFieldMask(req.UpdateMask),
		SubmitImmediately:     req.SubmitImmediately,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.ResubmitApplicationRes{Application: out.GetApplication()}, nil
}

func (c *ControllerV1) GetMyApplication(ctx context.Context, req *sellerv1.GetMyApplicationReq) (res *sellerv1.GetMyApplicationRes, err error) {
	out, err := c.svc.GetMyApplication(withRequestMetadata(ctx), &pb.GetMyApplicationReq{
		ApplicationNo: req.ApplicationNo,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.GetMyApplicationRes{Application: out.GetApplication()}, nil
}

func (c *ControllerV1) ListMyApplications(ctx context.Context, req *sellerv1.ListMyApplicationsReq) (res *sellerv1.ListMyApplicationsRes, err error) {
	out, err := c.svc.ListMyApplications(withRequestMetadata(ctx), &pb.ListMyApplicationsReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Statuses: toPBApplicationStatuses(req.Statuses),
	})
	if err != nil {
		return nil, err
	}
	return toHTTPListMyApplicationsRes(out), nil
}

func (c *ControllerV1) ListApplications(ctx context.Context, req *sellerv1.ListApplicationsReq) (res *sellerv1.ListApplicationsRes, err error) {
	out, err := c.svc.ListApplications(withRequestMetadata(ctx), &pb.ListApplicationsReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Statuses: toPBApplicationStatuses(req.Statuses),
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	return toHTTPListApplicationsRes(out), nil
}

func (c *ControllerV1) GetApplicationDetail(ctx context.Context, req *sellerv1.GetApplicationDetailReq) (res *sellerv1.GetApplicationDetailRes, err error) {
	out, err := c.svc.GetApplicationDetail(withRequestMetadata(ctx), &pb.GetApplicationDetailReq{
		ApplicationNo: req.ApplicationNo,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.GetApplicationDetailRes{Application: out.GetApplication()}, nil
}

func (c *ControllerV1) ApproveApplication(ctx context.Context, req *sellerv1.ApproveApplicationReq) (res *sellerv1.ApproveApplicationRes, err error) {
	out, err := c.svc.ApproveApplication(withRequestMetadata(ctx), &pb.ApproveApplicationReq{
		ApplicationNo:   req.ApplicationNo,
		ExpectedVersion: req.ExpectedVersion,
		ReviewComment:   req.ReviewComment,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.ApproveApplicationRes{
		ApplicationNo: out.GetApplicationNo(),
		ShopNo:        out.GetShopNo(),
		ShopStatus:    int32(out.GetShopStatus()),
	}, nil
}

func (c *ControllerV1) RejectApplication(ctx context.Context, req *sellerv1.RejectApplicationReq) (res *sellerv1.RejectApplicationRes, err error) {
	out, err := c.svc.RejectApplication(withRequestMetadata(ctx), &pb.RejectApplicationReq{
		ApplicationNo:    req.ApplicationNo,
		ExpectedVersion:  req.ExpectedVersion,
		RejectReasonCode: req.RejectReasonCode,
		RejectComment:    req.RejectComment,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.RejectApplicationRes{
		ApplicationNo: out.GetApplicationNo(),
		Status:        int32(out.GetStatus()),
	}, nil
}

func (c *ControllerV1) FreezeShop(ctx context.Context, req *sellerv1.FreezeShopReq) (res *sellerv1.FreezeShopRes, err error) {
	out, err := c.svc.FreezeShop(withRequestMetadata(ctx), &pb.FreezeShopReq{
		ShopNo:          req.ShopNo,
		ExpectedVersion: req.ExpectedVersion,
		ReasonCode:      req.ReasonCode,
		Reason:          req.Reason,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.FreezeShopRes{
		ShopNo:     out.GetShopNo(),
		ShopStatus: int32(out.GetShopStatus()),
	}, nil
}

func (c *ControllerV1) CloseShop(ctx context.Context, req *sellerv1.CloseShopReq) (res *sellerv1.CloseShopRes, err error) {
	out, err := c.svc.CloseShop(withRequestMetadata(ctx), &pb.CloseShopReq{
		ShopNo:          req.ShopNo,
		ExpectedVersion: req.ExpectedVersion,
		ReasonCode:      req.ReasonCode,
		Reason:          req.Reason,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.CloseShopRes{
		ShopNo:     out.GetShopNo(),
		ShopStatus: int32(out.GetShopStatus()),
	}, nil
}

func (c *ControllerV1) GetShopByNo(ctx context.Context, req *sellerv1.GetShopByNoReq) (res *sellerv1.GetShopByNoRes, err error) {
	out, err := c.svc.GetShopByNo(withRequestMetadata(ctx), &pb.GetShopByNoReq{
		ShopNo: req.ShopNo,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.GetShopByNoRes{Shop: out.GetShop()}, nil
}

func (c *ControllerV1) BatchGetShopsByNo(ctx context.Context, req *sellerv1.BatchGetShopsByNoReq) (res *sellerv1.BatchGetShopsByNoRes, err error) {
	out, err := c.svc.BatchGetShopsByNo(withRequestMetadata(ctx), &pb.BatchGetShopsByNoReq{
		ShopNos: req.ShopNos,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.BatchGetShopsByNoRes{Shops: out.GetShops()}, nil
}

func (c *ControllerV1) ListShopsByOwnerUserId(ctx context.Context, req *sellerv1.ListShopsByOwnerUserIdReq) (res *sellerv1.ListShopsByOwnerUserIdRes, err error) {
	out, err := c.svc.ListShopsByOwnerUserId(withRequestMetadata(ctx), &pb.ListShopsByOwnerUserIdReq{
		OwnerUserId: req.OwnerUserId,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.ListShopsByOwnerUserIdRes{Shops: out.GetShops()}, nil
}

func (c *ControllerV1) IsUserShopOwner(ctx context.Context, req *sellerv1.IsUserShopOwnerReq) (res *sellerv1.IsUserShopOwnerRes, err error) {
	out, err := c.svc.IsUserShopOwner(withRequestMetadata(ctx), &pb.IsUserShopOwnerReq{
		UserId: req.UserId,
		ShopNo: req.ShopNo,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.IsUserShopOwnerRes{IsOwner: out.GetIsOwner()}, nil
}

