package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// Controller binds gRPC handlers to seller shop service layer.
type Controller struct {
	v1.UnimplementedSellerApplicationServiceServer
	v1.UnimplementedAdminSellerServiceServer
	v1.UnimplementedInternalShopServiceServer
	svc service.ISellerShop
}

// Register 澶勭悊娉ㄥ唽涓绘祦绋嬪強鍒濆鍖栧姩浣溿€
func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{
		svc: service.SellerShop(),
	}
	v1.RegisterSellerApplicationServiceServer(s.Server, ctrl)
	v1.RegisterAdminSellerServiceServer(s.Server, ctrl)
	v1.RegisterInternalShopServiceServer(s.Server, ctrl)
}

// CreateApplicationDraft 鍒涘缓鏂拌褰曞苟杩斿洖鍒涘缓缁撴灉銆
func (c *Controller) CreateApplicationDraft(ctx context.Context, req *v1.CreateApplicationDraftReq) (*v1.CreateApplicationDraftRes, error) {
	return c.svc.CreateApplicationDraft(ctx, req)
}

// UpdateApplicationDraft 鎸夋潯浠舵洿鏂版暟鎹苟杩斿洖鏈€鏂扮粨鏋溿€
func (c *Controller) UpdateApplicationDraft(ctx context.Context, req *v1.UpdateApplicationDraftReq) (*v1.UpdateApplicationDraftRes, error) {
	return c.svc.UpdateApplicationDraft(ctx, req)
}

// SubmitApplication 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) SubmitApplication(ctx context.Context, req *v1.SubmitApplicationReq) (*v1.SubmitApplicationRes, error) {
	return c.svc.SubmitApplication(ctx, req)
}

// ResubmitApplication 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) ResubmitApplication(ctx context.Context, req *v1.ResubmitApplicationReq) (*v1.ResubmitApplicationRes, error) {
	return c.svc.ResubmitApplication(ctx, req)
}

// GetMyApplication 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (c *Controller) GetMyApplication(ctx context.Context, req *v1.GetMyApplicationReq) (*v1.GetMyApplicationRes, error) {
	return c.svc.GetMyApplication(ctx, req)
}

// ListMyApplications 鎸夋潯浠惰鍙栧苟杩斿洖鍒楄〃缁撴灉銆
func (c *Controller) ListMyApplications(ctx context.Context, req *v1.ListMyApplicationsReq) (*v1.ListMyApplicationsRes, error) {
	return c.svc.ListMyApplications(ctx, req)
}

// ListApplications 鎸夋潯浠惰鍙栧苟杩斿洖鍒楄〃缁撴灉銆
func (c *Controller) ListApplications(ctx context.Context, req *v1.ListApplicationsReq) (*v1.ListApplicationsRes, error) {
	return c.svc.ListApplications(ctx, req)
}

// GetApplicationDetail 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (c *Controller) GetApplicationDetail(ctx context.Context, req *v1.GetApplicationDetailReq) (*v1.GetApplicationDetailRes, error) {
	return c.svc.GetApplicationDetail(ctx, req)
}

// ApproveApplication 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) ApproveApplication(ctx context.Context, req *v1.ApproveApplicationReq) (*v1.ApproveApplicationRes, error) {
	return c.svc.ApproveApplication(ctx, req)
}

// RejectApplication 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) RejectApplication(ctx context.Context, req *v1.RejectApplicationReq) (*v1.RejectApplicationRes, error) {
	return c.svc.RejectApplication(ctx, req)
}

// FreezeShop 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) FreezeShop(ctx context.Context, req *v1.FreezeShopReq) (*v1.FreezeShopRes, error) {
	return c.svc.FreezeShop(ctx, req)
}

// CloseShop 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) CloseShop(ctx context.Context, req *v1.CloseShopReq) (*v1.CloseShopRes, error) {
	return c.svc.CloseShop(ctx, req)
}

// GetShopByNo 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (c *Controller) GetShopByNo(ctx context.Context, req *v1.GetShopByNoReq) (*v1.GetShopByNoRes, error) {
	return c.svc.GetShopByNo(ctx, req)
}

// BatchGetShopsByNo 鎵归噺澶勭悊璇锋眰锛屽噺灏戝線杩斿紑閿€銆
func (c *Controller) BatchGetShopsByNo(ctx context.Context, req *v1.BatchGetShopsByNoReq) (*v1.BatchGetShopsByNoRes, error) {
	return c.svc.BatchGetShopsByNo(ctx, req)
}

// ListShopsByOwnerUserId 鎸夋潯浠惰鍙栧苟杩斿洖鍒楄〃缁撴灉銆
func (c *Controller) ListShopsByOwnerUserId(ctx context.Context, req *v1.ListShopsByOwnerUserIdReq) (*v1.ListShopsByOwnerUserIdRes, error) {
	return c.svc.ListShopsByOwnerUserId(ctx, req)
}

// IsUserShopOwner 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) IsUserShopOwner(ctx context.Context, req *v1.IsUserShopOwnerReq) (*v1.IsUserShopOwnerRes, error) {
	return c.svc.IsUserShopOwner(ctx, req)
}
