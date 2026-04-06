// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 你可以删除这些注释，并按需手动维护该接口文件。
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
		ListStoreCategories(ctx context.Context, req *v1.ListStoreCategoriesReq) (*v1.ListStoreCategoriesRes, error)
		CreateStoreCategory(ctx context.Context, req *v1.CreateStoreCategoryReq) (*v1.CreateStoreCategoryRes, error)
		UpdateStoreCategory(ctx context.Context, req *v1.UpdateStoreCategoryReq) (*v1.UpdateStoreCategoryRes, error)
		SortStoreCategories(ctx context.Context, req *v1.SortStoreCategoriesReq) (*v1.SortStoreCategoriesRes, error)
		DeleteStoreCategory(ctx context.Context, req *v1.DeleteStoreCategoryReq) (*v1.DeleteStoreCategoryRes, error)
		GetProductStoreCategoryBinding(ctx context.Context, req *v1.GetProductStoreCategoryBindingReq) (*v1.GetProductStoreCategoryBindingRes, error)
		BatchGetProductStoreCategoryBindings(ctx context.Context, req *v1.BatchGetProductStoreCategoryBindingsReq) (*v1.BatchGetProductStoreCategoryBindingsRes, error)
		UpdateProductStoreCategoryBinding(ctx context.Context, req *v1.UpdateProductStoreCategoryBindingReq) (*v1.UpdateProductStoreCategoryBindingRes, error)
		ListBuyerStoreCategories(ctx context.Context, req *v1.ListBuyerStoreCategoriesReq) (*v1.ListBuyerStoreCategoriesRes, error)
		ReconcileStoreCategoryCounts(ctx context.Context, req *v1.ReconcileStoreCategoryCountsReq) (*v1.ReconcileStoreCategoryCountsRes, error)
	}
)

var (
	localSellerShop ISellerShop
)

// SellerShop 返回已注册的卖家店铺服务实现。
func SellerShop() ISellerShop {
	if localSellerShop == nil {
		panic("implement not found for interface ISellerShop, forgot register?")
	}
	return localSellerShop
}

// RegisterSellerShop 注册卖家店铺服务实现，通常在 internal/logic 的 init 中调用。
func RegisterSellerShop(i ISellerShop) {
	localSellerShop = i
}
