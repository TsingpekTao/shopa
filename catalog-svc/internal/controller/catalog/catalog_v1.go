package catalog

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/catalog-svc/api/catalog/v1"
	pb "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *ControllerV1) CreateProductDraft(ctx context.Context, req *v1.CreateProductDraftReq) (*v1.CreateProductDraftRes, error) {
	return c.catalog.CreateProductDraft(ctx, &req.CreateProductDraftReq)
}

func (c *ControllerV1) UpdateProductDraft(ctx context.Context, req *v1.UpdateProductDraftReq) (*v1.UpdateProductDraftRes, error) {
	return c.catalog.UpdateProductDraft(ctx, &req.UpdateProductDraftReq)
}

func (c *ControllerV1) UpsertSkuDrafts(ctx context.Context, req *v1.UpsertSkuDraftsReq) (*v1.UpsertSkuDraftsRes, error) {
	return c.catalog.UpsertSkuDrafts(ctx, &req.UpsertSkuDraftsReq)
}

func (c *ControllerV1) SubmitProductReview(ctx context.Context, req *v1.SubmitProductReviewReq) (*v1.SubmitProductReviewRes, error) {
	return c.catalog.SubmitProductReview(ctx, &req.SubmitProductReviewReq)
}

func (c *ControllerV1) ResubmitProductReview(ctx context.Context, req *v1.ResubmitProductReviewReq) (*v1.ResubmitProductReviewRes, error) {
	return c.catalog.ResubmitProductReview(ctx, &req.ResubmitProductReviewReq)
}

func (c *ControllerV1) SetProductOnShelf(ctx context.Context, req *v1.SetProductOnShelfReq) (*v1.SetProductOnShelfRes, error) {
	return c.catalog.SetProductOnShelf(ctx, &req.SetProductOnShelfReq)
}

func (c *ControllerV1) SetProductOffShelf(ctx context.Context, req *v1.SetProductOffShelfReq) (*v1.SetProductOffShelfRes, error) {
	return c.catalog.SetProductOffShelf(ctx, &req.SetProductOffShelfReq)
}

func (c *ControllerV1) DeleteProductDraft(ctx context.Context, req *v1.DeleteProductDraftReq) (*v1.DeleteProductDraftRes, error) {
	_, err := c.catalog.DeleteProductDraft(ctx, &req.DeleteProductDraftReq)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteProductDraftRes{Success: true}, nil
}

func (c *ControllerV1) GetMyProduct(ctx context.Context, req *v1.GetMyProductReq) (*v1.GetMyProductRes, error) {
	return c.catalog.GetMyProduct(ctx, &pb.GetMyProductReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ListMyProducts(ctx context.Context, req *v1.ListMyProductsReq) (*v1.ListMyProductsRes, error) {
	return c.catalog.ListMyProducts(ctx, &pb.ListMyProductsReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Statuses: req.Statuses,
		Keyword:  req.Keyword,
	})
}

func (c *ControllerV1) ListReviewTasks(ctx context.Context, req *v1.ListReviewTasksReq) (*v1.ListReviewTasksRes, error) {
	return c.catalog.ListReviewTasks(ctx, &pb.ListReviewTasksReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Statuses: req.Statuses,
		Keyword:  req.Keyword,
	})
}

func (c *ControllerV1) GetReviewDetail(ctx context.Context, req *v1.GetReviewDetailReq) (*v1.GetReviewDetailRes, error) {
	return c.catalog.GetReviewDetail(ctx, &pb.GetReviewDetailReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ApproveProduct(ctx context.Context, req *v1.ApproveProductReq) (*v1.ApproveProductRes, error) {
	return c.catalog.ApproveProduct(ctx, &req.ApproveProductReq)
}

func (c *ControllerV1) RejectProduct(ctx context.Context, req *v1.RejectProductReq) (*v1.RejectProductRes, error) {
	return c.catalog.RejectProduct(ctx, &req.RejectProductReq)
}

func (c *ControllerV1) FreezeProduct(ctx context.Context, req *v1.FreezeProductReq) (*v1.FreezeProductRes, error) {
	return c.catalog.FreezeProduct(ctx, &req.FreezeProductReq)
}

func (c *ControllerV1) UnfreezeProduct(ctx context.Context, req *v1.UnfreezeProductReq) (*v1.UnfreezeProductRes, error) {
	return c.catalog.UnfreezeProduct(ctx, &req.UnfreezeProductReq)
}

func (c *ControllerV1) ForceOffShelf(ctx context.Context, req *v1.ForceOffShelfReq) (*v1.ForceOffShelfRes, error) {
	return c.catalog.ForceOffShelf(ctx, &req.ForceOffShelfReq)
}

func (c *ControllerV1) GetProductDetail(ctx context.Context, req *v1.GetProductDetailReq) (*v1.GetProductDetailRes, error) {
	return c.catalog.GetProductDetail(ctx, &pb.GetProductDetailReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ListProducts(ctx context.Context, req *v1.ListProductsReq) (*v1.ListProductsRes, error) {
	return c.catalog.ListProducts(ctx, &pb.ListProductsReq{
		CategoryId: req.CategoryId,
		Page:       req.Page,
		PageSize:   req.PageSize,
		SortBy:     req.SortBy,
	})
}

func (c *ControllerV1) SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error) {
	return c.catalog.SearchProducts(ctx, &pb.SearchProductsReq{
		Keyword:    req.Keyword,
		CategoryId: req.CategoryId,
		Page:       req.Page,
		PageSize:   req.PageSize,
		SortBy:     req.SortBy,
	})
}

var _ = emptypb.Empty{}
