package catalog

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/catalog-svc/api/catalog/v1"
	pb "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *ControllerV1) CreateProductDraft(ctx context.Context, req *pb.CreateProductDraftReq) (*pb.CreateProductDraftRes, error) {
	return c.catalog.CreateProductDraft(ctx, req)
}

func (c *ControllerV1) UpdateProductDraft(ctx context.Context, req *pb.UpdateProductDraftReq) (*pb.UpdateProductDraftRes, error) {
	return c.catalog.UpdateProductDraft(ctx, req)
}

func (c *ControllerV1) UpsertSkuDrafts(ctx context.Context, req *pb.UpsertSkuDraftsReq) (*pb.UpsertSkuDraftsRes, error) {
	return c.catalog.UpsertSkuDrafts(ctx, req)
}

func (c *ControllerV1) SubmitProductReview(ctx context.Context, req *pb.SubmitProductReviewReq) (*pb.SubmitProductReviewRes, error) {
	return c.catalog.SubmitProductReview(ctx, req)
}

func (c *ControllerV1) ResubmitProductReview(ctx context.Context, req *pb.ResubmitProductReviewReq) (*pb.ResubmitProductReviewRes, error) {
	return c.catalog.ResubmitProductReview(ctx, req)
}

func (c *ControllerV1) SetProductOnShelf(ctx context.Context, req *pb.SetProductOnShelfReq) (*pb.SetProductOnShelfRes, error) {
	return c.catalog.SetProductOnShelf(ctx, req)
}

func (c *ControllerV1) SetProductOffShelf(ctx context.Context, req *pb.SetProductOffShelfReq) (*pb.SetProductOffShelfRes, error) {
	return c.catalog.SetProductOffShelf(ctx, req)
}

func (c *ControllerV1) DeleteProductDraft(ctx context.Context, req *pb.DeleteProductDraftReq) (*emptypb.Empty, error) {
	_, err := c.catalog.DeleteProductDraft(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c *ControllerV1) GetMyProduct(ctx context.Context, req *pb.GetMyProductReq) (*pb.GetMyProductRes, error) {
	return c.catalog.GetMyProduct(ctx, &pb.GetMyProductReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ListMyProducts(ctx context.Context, req *pb.ListMyProductsReq) (*pb.ListMyProductsRes, error) {
	return c.catalog.ListMyProducts(ctx, &pb.ListMyProductsReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Statuses: req.Statuses,
		Keyword:  req.Keyword,
	})
}

func (c *ControllerV1) ListReviewTasks(ctx context.Context, req *pb.ListReviewTasksReq) (*pb.ListReviewTasksRes, error) {
	return c.catalog.ListReviewTasks(ctx, &pb.ListReviewTasksReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Statuses: req.Statuses,
		Keyword:  req.Keyword,
	})
}

func (c *ControllerV1) GetReviewDetail(ctx context.Context, req *pb.GetReviewDetailReq) (*pb.GetReviewDetailRes, error) {
	return c.catalog.GetReviewDetail(ctx, &pb.GetReviewDetailReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ApproveProduct(ctx context.Context, req *pb.ApproveProductReq) (*pb.ApproveProductRes, error) {
	return c.catalog.ApproveProduct(ctx, req)
}

func (c *ControllerV1) RejectProduct(ctx context.Context, req *pb.RejectProductReq) (*pb.RejectProductRes, error) {
	return c.catalog.RejectProduct(ctx, req)
}

func (c *ControllerV1) FreezeProduct(ctx context.Context, req *pb.FreezeProductReq) (*pb.FreezeProductRes, error) {
	return c.catalog.FreezeProduct(ctx, req)
}

func (c *ControllerV1) UnfreezeProduct(ctx context.Context, req *pb.UnfreezeProductReq) (*pb.UnfreezeProductRes, error) {
	return c.catalog.UnfreezeProduct(ctx, req)
}

func (c *ControllerV1) ForceOffShelf(ctx context.Context, req *pb.ForceOffShelfReq) (*pb.ForceOffShelfRes, error) {
	return c.catalog.ForceOffShelf(ctx, req)
}

func (c *ControllerV1) GetProductDetail(ctx context.Context, req *httpv1.GetProductDetailReq) (*httpv1.GetProductDetailRes, error) {
	return c.catalog.GetProductDetail(ctx, &pb.GetProductDetailReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ListProducts(ctx context.Context, req *httpv1.ListProductsReq) (*httpv1.ListProductsRes, error) {
	return c.catalog.ListProducts(ctx, &pb.ListProductsReq{
		CategoryId: req.CategoryId,
		Page:       req.Page,
		PageSize:   req.PageSize,
		SortBy:     req.SortBy,
	})
}

func (c *ControllerV1) SearchProducts(ctx context.Context, req *httpv1.SearchProductsReq) (*httpv1.SearchProductsRes, error) {
	return c.catalog.SearchProducts(ctx, &pb.SearchProductsReq{
		Keyword:    req.Keyword,
		CategoryId: req.CategoryId,
		Page:       req.Page,
		PageSize:   req.PageSize,
		SortBy:     req.SortBy,
	})
}

var _ = emptypb.Empty{}

