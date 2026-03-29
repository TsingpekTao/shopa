package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedSellerApplicationServiceServer
	v1.UnimplementedAdminSellerServiceServer
	v1.UnimplementedInternalShopServiceServer
	svc service.ISellerShop
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{svc: service.SellerShop()}
	v1.RegisterSellerApplicationServiceServer(s.Server, ctrl)
	v1.RegisterAdminSellerServiceServer(s.Server, ctrl)
	v1.RegisterInternalShopServiceServer(s.Server, ctrl)
}

func (c *Controller) CreateApplicationDraft(ctx context.Context, req *v1.CreateApplicationDraftReq) (*v1.CreateApplicationDraftRes, error) {
	return c.svc.CreateApplicationDraft(ctx, req)
}

func (c *Controller) UpdateApplicationDraft(ctx context.Context, req *v1.UpdateApplicationDraftReq) (*v1.UpdateApplicationDraftRes, error) {
	return c.svc.UpdateApplicationDraft(ctx, req)
}

func (c *Controller) SubmitApplication(ctx context.Context, req *v1.SubmitApplicationReq) (*v1.SubmitApplicationRes, error) {
	return c.svc.SubmitApplication(ctx, req)
}

func (c *Controller) ResubmitApplication(ctx context.Context, req *v1.ResubmitApplicationReq) (*v1.ResubmitApplicationRes, error) {
	return c.svc.ResubmitApplication(ctx, req)
}

func (c *Controller) GetMyApplication(ctx context.Context, req *v1.GetMyApplicationReq) (*v1.GetMyApplicationRes, error) {
	return c.svc.GetMyApplication(ctx, req)
}

func (c *Controller) ListMyApplications(ctx context.Context, req *v1.ListMyApplicationsReq) (*v1.ListMyApplicationsRes, error) {
	return c.svc.ListMyApplications(ctx, req)
}

func (c *Controller) ListApplications(ctx context.Context, req *v1.ListApplicationsReq) (*v1.ListApplicationsRes, error) {
	return c.svc.ListApplications(ctx, req)
}

func (c *Controller) GetApplicationDetail(ctx context.Context, req *v1.GetApplicationDetailReq) (*v1.GetApplicationDetailRes, error) {
	return c.svc.GetApplicationDetail(ctx, req)
}

func (c *Controller) ApproveApplication(ctx context.Context, req *v1.ApproveApplicationReq) (*v1.ApproveApplicationRes, error) {
	return c.svc.ApproveApplication(ctx, req)
}

func (c *Controller) RejectApplication(ctx context.Context, req *v1.RejectApplicationReq) (*v1.RejectApplicationRes, error) {
	return c.svc.RejectApplication(ctx, req)
}

func (c *Controller) FreezeShop(ctx context.Context, req *v1.FreezeShopReq) (*v1.FreezeShopRes, error) {
	return c.svc.FreezeShop(ctx, req)
}

func (c *Controller) CloseShop(ctx context.Context, req *v1.CloseShopReq) (*v1.CloseShopRes, error) {
	return c.svc.CloseShop(ctx, req)
}

func (c *Controller) GetShopByNo(ctx context.Context, req *v1.GetShopByNoReq) (*v1.GetShopByNoRes, error) {
	return c.svc.GetShopByNo(ctx, req)
}

func (c *Controller) BatchGetShopsByNo(ctx context.Context, req *v1.BatchGetShopsByNoReq) (*v1.BatchGetShopsByNoRes, error) {
	return c.svc.BatchGetShopsByNo(ctx, req)
}

func (c *Controller) ListShopsByOwnerUserId(ctx context.Context, req *v1.ListShopsByOwnerUserIdReq) (*v1.ListShopsByOwnerUserIdRes, error) {
	return c.svc.ListShopsByOwnerUserId(ctx, req)
}

func (c *Controller) IsUserShopOwner(ctx context.Context, req *v1.IsUserShopOwnerReq) (*v1.IsUserShopOwnerRes, error) {
	return c.svc.IsUserShopOwner(ctx, req)
}
