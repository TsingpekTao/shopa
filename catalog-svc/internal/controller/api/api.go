package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"github.com/TsingpekTao/shopa/catalog-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	v1.UnimplementedSellerProductServiceServer
	v1.UnimplementedAdminProductReviewServiceServer
	v1.UnimplementedBuyerCatalogServiceServer
	v1.UnimplementedInternalCatalogServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterSellerProductServiceServer(s.Server, ctrl)
	v1.RegisterAdminProductReviewServiceServer(s.Server, ctrl)
	v1.RegisterBuyerCatalogServiceServer(s.Server, ctrl)
	v1.RegisterInternalCatalogServiceServer(s.Server, ctrl)
}

func (*Controller) CreateProductDraft(ctx context.Context, req *v1.CreateProductDraftReq) (*v1.CreateProductDraftRes, error) {
	return service.Catalog().CreateProductDraft(ctx, req)
}

func (*Controller) UpdateProductDraft(ctx context.Context, req *v1.UpdateProductDraftReq) (*v1.UpdateProductDraftRes, error) {
	return service.Catalog().UpdateProductDraft(ctx, req)
}

func (*Controller) UpsertSkuDrafts(ctx context.Context, req *v1.UpsertSkuDraftsReq) (*v1.UpsertSkuDraftsRes, error) {
	return service.Catalog().UpsertSkuDrafts(ctx, req)
}

func (*Controller) SubmitProductReview(ctx context.Context, req *v1.SubmitProductReviewReq) (*v1.SubmitProductReviewRes, error) {
	return service.Catalog().SubmitProductReview(ctx, req)
}

func (*Controller) ResubmitProductReview(ctx context.Context, req *v1.ResubmitProductReviewReq) (*v1.ResubmitProductReviewRes, error) {
	return service.Catalog().ResubmitProductReview(ctx, req)
}

func (*Controller) SetProductOnShelf(ctx context.Context, req *v1.SetProductOnShelfReq) (*v1.SetProductOnShelfRes, error) {
	return service.Catalog().SetProductOnShelf(ctx, req)
}

func (*Controller) SetProductOffShelf(ctx context.Context, req *v1.SetProductOffShelfReq) (*v1.SetProductOffShelfRes, error) {
	return service.Catalog().SetProductOffShelf(ctx, req)
}

func (*Controller) DeleteProductDraft(ctx context.Context, req *v1.DeleteProductDraftReq) (*emptypb.Empty, error) {
	return service.Catalog().DeleteProductDraft(ctx, req)
}

func (*Controller) GetMyProduct(ctx context.Context, req *v1.GetMyProductReq) (*v1.GetMyProductRes, error) {
	return service.Catalog().GetMyProduct(ctx, req)
}

func (*Controller) ListMyProducts(ctx context.Context, req *v1.ListMyProductsReq) (*v1.ListMyProductsRes, error) {
	return service.Catalog().ListMyProducts(ctx, req)
}

func (*Controller) ListReviewTasks(ctx context.Context, req *v1.ListReviewTasksReq) (*v1.ListReviewTasksRes, error) {
	return service.Catalog().ListReviewTasks(ctx, req)
}

func (*Controller) GetReviewDetail(ctx context.Context, req *v1.GetReviewDetailReq) (*v1.GetReviewDetailRes, error) {
	return service.Catalog().GetReviewDetail(ctx, req)
}

func (*Controller) ApproveProduct(ctx context.Context, req *v1.ApproveProductReq) (*v1.ApproveProductRes, error) {
	return service.Catalog().ApproveProduct(ctx, req)
}

func (*Controller) RejectProduct(ctx context.Context, req *v1.RejectProductReq) (*v1.RejectProductRes, error) {
	return service.Catalog().RejectProduct(ctx, req)
}

func (*Controller) FreezeProduct(ctx context.Context, req *v1.FreezeProductReq) (*v1.FreezeProductRes, error) {
	return service.Catalog().FreezeProduct(ctx, req)
}

func (*Controller) UnfreezeProduct(ctx context.Context, req *v1.UnfreezeProductReq) (*v1.UnfreezeProductRes, error) {
	return service.Catalog().UnfreezeProduct(ctx, req)
}

func (*Controller) ForceOffShelf(ctx context.Context, req *v1.ForceOffShelfReq) (*v1.ForceOffShelfRes, error) {
	return service.Catalog().ForceOffShelf(ctx, req)
}

func (*Controller) GetProductDetail(ctx context.Context, req *v1.GetProductDetailReq) (*v1.GetProductDetailRes, error) {
	return service.Catalog().GetProductDetail(ctx, req)
}

func (*Controller) ListProducts(ctx context.Context, req *v1.ListProductsReq) (*v1.ListProductsRes, error) {
	return service.Catalog().ListProducts(ctx, req)
}

func (*Controller) SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error) {
	return service.Catalog().SearchProducts(ctx, req)
}

func (*Controller) BatchGetSpuByNo(ctx context.Context, req *v1.BatchGetSpuByNoReq) (*v1.BatchGetSpuByNoRes, error) {
	return service.Catalog().BatchGetSpuByNo(ctx, req)
}

func (*Controller) BatchGetSkuByNo(ctx context.Context, req *v1.BatchGetSkuByNoReq) (*v1.BatchGetSkuByNoRes, error) {
	return service.Catalog().BatchGetSkuByNo(ctx, req)
}

func (*Controller) GetSkuSnapshotForOrder(ctx context.Context, req *v1.GetSkuSnapshotForOrderReq) (*v1.GetSkuSnapshotForOrderRes, error) {
	return service.Catalog().GetSkuSnapshotForOrder(ctx, req)
}

func (*Controller) UpsertSkuStockProjection(ctx context.Context, req *v1.UpsertSkuStockProjectionReq) (*v1.UpsertSkuStockProjectionRes, error) {
	return service.Catalog().UpsertSkuStockProjection(ctx, req)
}

func (*Controller) RecomputeSpuAggregation(ctx context.Context, req *v1.RecomputeSpuAggregationReq) (*v1.RecomputeSpuAggregationRes, error) {
	return service.Catalog().RecomputeSpuAggregation(ctx, req)
}
