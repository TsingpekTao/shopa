package catalog

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/catalog-svc/api/catalog/v1"
)

// ICatalogV1 defines HTTP handlers for catalog business APIs.
type ICatalogV1 interface {
	CreateProductDraft(ctx context.Context, req *v1.CreateProductDraftReq) (res *v1.CreateProductDraftRes, err error)
	UpdateProductDraft(ctx context.Context, req *v1.UpdateProductDraftReq) (res *v1.UpdateProductDraftRes, err error)
	UpsertSkuDrafts(ctx context.Context, req *v1.UpsertSkuDraftsReq) (res *v1.UpsertSkuDraftsRes, err error)
	SubmitProductReview(ctx context.Context, req *v1.SubmitProductReviewReq) (res *v1.SubmitProductReviewRes, err error)
	ResubmitProductReview(ctx context.Context, req *v1.ResubmitProductReviewReq) (res *v1.ResubmitProductReviewRes, err error)
	SetProductOnShelf(ctx context.Context, req *v1.SetProductOnShelfReq) (res *v1.SetProductOnShelfRes, err error)
	SetProductOffShelf(ctx context.Context, req *v1.SetProductOffShelfReq) (res *v1.SetProductOffShelfRes, err error)
	DeleteProductDraft(ctx context.Context, req *v1.DeleteProductDraftReq) (res *v1.DeleteProductDraftRes, err error)
	GetMyProduct(ctx context.Context, req *v1.GetMyProductReq) (res *v1.GetMyProductRes, err error)
	ListMyProducts(ctx context.Context, req *v1.ListMyProductsReq) (res *v1.ListMyProductsRes, err error)

	ListReviewTasks(ctx context.Context, req *v1.ListReviewTasksReq) (res *v1.ListReviewTasksRes, err error)
	GetReviewDetail(ctx context.Context, req *v1.GetReviewDetailReq) (res *v1.GetReviewDetailRes, err error)
	ApproveProduct(ctx context.Context, req *v1.ApproveProductReq) (res *v1.ApproveProductRes, err error)
	RejectProduct(ctx context.Context, req *v1.RejectProductReq) (res *v1.RejectProductRes, err error)
	FreezeProduct(ctx context.Context, req *v1.FreezeProductReq) (res *v1.FreezeProductRes, err error)
	UnfreezeProduct(ctx context.Context, req *v1.UnfreezeProductReq) (res *v1.UnfreezeProductRes, err error)
	ForceOffShelf(ctx context.Context, req *v1.ForceOffShelfReq) (res *v1.ForceOffShelfRes, err error)

	GetProductDetail(ctx context.Context, req *v1.GetProductDetailReq) (res *v1.GetProductDetailRes, err error)
	ListProducts(ctx context.Context, req *v1.ListProductsReq) (res *v1.ListProductsRes, err error)
	SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (res *v1.SearchProductsRes, err error)
}
