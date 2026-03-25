// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
)

type (
	ISellerShop interface {
		CreateApplicationDraft(ctx context.Context, req *v1.CreateApplicationDraftReq) (*v1.CreateApplicationDraftRes, error)
		UpdateApplicationDraft(ctx context.Context, req *v1.UpdateApplicationDraftReq) (*v1.UpdateApplicationDraftRes, error)
		SubmitApplication(ctx context.Context, req *v1.SubmitApplicationReq) (*v1.SubmitApplicationRes, error)
		ResubmitApplication(ctx context.Context, req *v1.ResubmitApplicationReq) (*v1.ResubmitApplicationRes, error)
		GetMyApplication(ctx context.Context, req *v1.GetMyApplicationReq) (*v1.GetMyApplicationRes, error)
		ListMyApplications(ctx context.Context, req *v1.ListMyApplicationsReq) (*v1.ListMyApplicationsRes, error)
		ListApplications(ctx context.Context, req *v1.ListApplicationsReq) (*v1.ListApplicationsRes, error)
		GetApplicationDetail(ctx context.Context, req *v1.GetApplicationDetailReq) (*v1.GetApplicationDetailRes, error)
		ApproveApplication(ctx context.Context, req *v1.ApproveApplicationReq) (*v1.ApproveApplicationRes, error)
		RejectApplication(ctx context.Context, req *v1.RejectApplicationReq) (*v1.RejectApplicationRes, error)
		FreezeShop(ctx context.Context, req *v1.FreezeShopReq) (*v1.FreezeShopRes, error)
		CloseShop(ctx context.Context, req *v1.CloseShopReq) (*v1.CloseShopRes, error)
		GetShopByNo(ctx context.Context, req *v1.GetShopByNoReq) (*v1.GetShopByNoRes, error)
		BatchGetShopsByNo(ctx context.Context, req *v1.BatchGetShopsByNoReq) (*v1.BatchGetShopsByNoRes, error)
		ListShopsByOwnerUserId(ctx context.Context, req *v1.ListShopsByOwnerUserIdReq) (*v1.ListShopsByOwnerUserIdRes, error)
		IsUserShopOwner(ctx context.Context, req *v1.IsUserShopOwnerReq) (*v1.IsUserShopOwnerRes, error)
	}
)

var (
	localSellerShop ISellerShop
)

// SellerShop 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func SellerShop() ISellerShop {
	if localSellerShop == nil {
		panic("implement not found for interface ISellerShop, forgot register?")
	}
	return localSellerShop
}

// RegisterSellerShop 澶勭悊娉ㄥ唽涓绘祦绋嬪強鍒濆鍖栧姩浣溿€
func RegisterSellerShop(i ISellerShop) {
	localSellerShop = i
}
