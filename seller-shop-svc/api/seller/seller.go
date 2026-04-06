package seller

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/seller/v1"
)

// ISellerV1 瀹氫箟 seller-shop-svc 瀵瑰 HTTP API銆?
type ISellerV1 interface {
	// 鍗栧鐢宠
	CreateApplicationDraft(ctx context.Context, req *v1.CreateApplicationDraftReq) (res *v1.CreateApplicationDraftRes, err error)
	UpdateApplicationDraft(ctx context.Context, req *v1.UpdateApplicationDraftReq) (res *v1.UpdateApplicationDraftRes, err error)
	SubmitApplication(ctx context.Context, req *v1.SubmitApplicationReq) (res *v1.SubmitApplicationRes, err error)
	ResubmitApplication(ctx context.Context, req *v1.ResubmitApplicationReq) (res *v1.ResubmitApplicationRes, err error)
	GetMyApplication(ctx context.Context, req *v1.GetMyApplicationReq) (res *v1.GetMyApplicationRes, err error)
	ListMyApplications(ctx context.Context, req *v1.ListMyApplicationsReq) (res *v1.ListMyApplicationsRes, err error)

	// 绠＄悊绔?
	ListApplications(ctx context.Context, req *v1.ListApplicationsReq) (res *v1.ListApplicationsRes, err error)
	GetApplicationDetail(ctx context.Context, req *v1.GetApplicationDetailReq) (res *v1.GetApplicationDetailRes, err error)
	ApproveApplication(ctx context.Context, req *v1.ApproveApplicationReq) (res *v1.ApproveApplicationRes, err error)
	RejectApplication(ctx context.Context, req *v1.RejectApplicationReq) (res *v1.RejectApplicationRes, err error)
	FreezeShop(ctx context.Context, req *v1.FreezeShopReq) (res *v1.FreezeShopRes, err error)
	CloseShop(ctx context.Context, req *v1.CloseShopReq) (res *v1.CloseShopRes, err error)

	// 鍐呴儴鏌ヨ
	GetShopByNo(ctx context.Context, req *v1.GetShopByNoReq) (res *v1.GetShopByNoRes, err error)
	BatchGetShopsByNo(ctx context.Context, req *v1.BatchGetShopsByNoReq) (res *v1.BatchGetShopsByNoRes, err error)
	ListShopsByOwnerUserId(ctx context.Context, req *v1.ListShopsByOwnerUserIdReq) (res *v1.ListShopsByOwnerUserIdRes, err error)
	IsUserShopOwner(ctx context.Context, req *v1.IsUserShopOwnerReq) (res *v1.IsUserShopOwnerRes, err error)
	ListStoreCategories(ctx context.Context, req *v1.ListStoreCategoriesReq) (res *v1.ListStoreCategoriesRes, err error)
	CreateStoreCategory(ctx context.Context, req *v1.CreateStoreCategoryReq) (res *v1.CreateStoreCategoryRes, err error)
	UpdateStoreCategory(ctx context.Context, req *v1.UpdateStoreCategoryReq) (res *v1.UpdateStoreCategoryRes, err error)
	SortStoreCategories(ctx context.Context, req *v1.SortStoreCategoriesReq) (res *v1.SortStoreCategoriesRes, err error)
	DeleteStoreCategory(ctx context.Context, req *v1.DeleteStoreCategoryReq) (res *v1.DeleteStoreCategoryRes, err error)
	GetProductStoreCategoryBinding(ctx context.Context, req *v1.GetProductStoreCategoryBindingReq) (res *v1.GetProductStoreCategoryBindingRes, err error)
	BatchGetProductStoreCategoryBindings(ctx context.Context, req *v1.BatchGetProductStoreCategoryBindingsReq) (res *v1.BatchGetProductStoreCategoryBindingsRes, err error)
	UpdateProductStoreCategoryBinding(ctx context.Context, req *v1.UpdateProductStoreCategoryBindingReq) (res *v1.UpdateProductStoreCategoryBindingRes, err error)
	ListBuyerStoreCategories(ctx context.Context, req *v1.ListBuyerStoreCategoriesReq) (res *v1.ListBuyerStoreCategoriesRes, err error)
}
