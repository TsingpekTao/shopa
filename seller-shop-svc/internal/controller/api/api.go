package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// Controller 负责将 gRPC 处理器绑定到 seller shop 服务层。
type Controller struct {
	v1.UnimplementedSellerApplicationServiceServer
	v1.UnimplementedAdminSellerServiceServer
	v1.UnimplementedInternalShopServiceServer
	svc service.ISellerShop
}

// Register 将 seller shop 控制器注册到 gRPC 服务器。
func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{
		svc: service.SellerShop(),
	}
	v1.RegisterSellerApplicationServiceServer(s.Server, ctrl)
	v1.RegisterAdminSellerServiceServer(s.Server, ctrl)
	v1.RegisterInternalShopServiceServer(s.Server, ctrl)
}

// CreateApplicationDraft 创建店铺申请草稿。
func (c *Controller) CreateApplicationDraft(ctx context.Context, req *v1.CreateApplicationDraftReq) (*v1.CreateApplicationDraftRes, error) {
	return c.svc.CreateApplicationDraft(ctx, req)
}

// UpdateApplicationDraft 更新店铺申请草稿。
func (c *Controller) UpdateApplicationDraft(ctx context.Context, req *v1.UpdateApplicationDraftReq) (*v1.UpdateApplicationDraftRes, error) {
	return c.svc.UpdateApplicationDraft(ctx, req)
}

// SubmitApplication 提交店铺申请。
func (c *Controller) SubmitApplication(ctx context.Context, req *v1.SubmitApplicationReq) (*v1.SubmitApplicationRes, error) {
	return c.svc.SubmitApplication(ctx, req)
}

// ResubmitApplication 重新提交店铺申请。
func (c *Controller) ResubmitApplication(ctx context.Context, req *v1.ResubmitApplicationReq) (*v1.ResubmitApplicationRes, error) {
	return c.svc.ResubmitApplication(ctx, req)
}

// GetMyApplication 查询当前用户的申请。
func (c *Controller) GetMyApplication(ctx context.Context, req *v1.GetMyApplicationReq) (*v1.GetMyApplicationRes, error) {
	return c.svc.GetMyApplication(ctx, req)
}

// ListMyApplications 列出当前用户的申请列表。
func (c *Controller) ListMyApplications(ctx context.Context, req *v1.ListMyApplicationsReq) (*v1.ListMyApplicationsRes, error) {
	return c.svc.ListMyApplications(ctx, req)
}

// ListApplications 列出所有店铺申请。
func (c *Controller) ListApplications(ctx context.Context, req *v1.ListApplicationsReq) (*v1.ListApplicationsRes, error) {
	return c.svc.ListApplications(ctx, req)
}

// GetApplicationDetail 获取申请详情。
func (c *Controller) GetApplicationDetail(ctx context.Context, req *v1.GetApplicationDetailReq) (*v1.GetApplicationDetailRes, error) {
	return c.svc.GetApplicationDetail(ctx, req)
}

// ApproveApplication 审核并通过申请。
func (c *Controller) ApproveApplication(ctx context.Context, req *v1.ApproveApplicationReq) (*v1.ApproveApplicationRes, error) {
	return c.svc.ApproveApplication(ctx, req)
}

// RejectApplication 审核并拒绝申请。
func (c *Controller) RejectApplication(ctx context.Context, req *v1.RejectApplicationReq) (*v1.RejectApplicationRes, error) {
	return c.svc.RejectApplication(ctx, req)
}

// FreezeShop 冻结店铺。
func (c *Controller) FreezeShop(ctx context.Context, req *v1.FreezeShopReq) (*v1.FreezeShopRes, error) {
	return c.svc.FreezeShop(ctx, req)
}

// CloseShop 关闭店铺。
func (c *Controller) CloseShop(ctx context.Context, req *v1.CloseShopReq) (*v1.CloseShopRes, error) {
	return c.svc.CloseShop(ctx, req)
}

// GetShopByNo 根据店铺编号获取店铺。
func (c *Controller) GetShopByNo(ctx context.Context, req *v1.GetShopByNoReq) (*v1.GetShopByNoRes, error) {
	return c.svc.GetShopByNo(ctx, req)
}

// BatchGetShopsByNo 批量根据店铺编号获取店铺。
func (c *Controller) BatchGetShopsByNo(ctx context.Context, req *v1.BatchGetShopsByNoReq) (*v1.BatchGetShopsByNoRes, error) {
	return c.svc.BatchGetShopsByNo(ctx, req)
}

// ListShopsByOwnerUserId 列出指定所有者的店铺。
func (c *Controller) ListShopsByOwnerUserId(ctx context.Context, req *v1.ListShopsByOwnerUserIdReq) (*v1.ListShopsByOwnerUserIdRes, error) {
	return c.svc.ListShopsByOwnerUserId(ctx, req)
}

// IsUserShopOwner 判断用户是否为店铺所有者。
func (c *Controller) IsUserShopOwner(ctx context.Context, req *v1.IsUserShopOwnerReq) (*v1.IsUserShopOwnerRes, error) {
	return c.svc.IsUserShopOwner(ctx, req)
}


