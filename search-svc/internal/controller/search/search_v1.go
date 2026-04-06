package search

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/search-svc/api/search/v1"
	pb "github.com/TsingpekTao/shopa/search-svc/api/v1"
)

func (c *ControllerV1) SearchProducts(ctx context.Context, req *httpv1.SearchProductsReq) (*httpv1.SearchProductsRes, error) {
	return c.search.SearchProducts(ctx, &pb.SearchProductsReq{
		Query:              req.Query,
		ShopNo:             req.ShopNo,
		CategoryNo:         req.CategoryNo,
		StoreCategoryId:    req.StoreCategoryId,
		StoreCategoryLevel: req.StoreCategoryLevel,
		SortCode:           req.SortCode,
		PageSize:           req.PageSize,
		NextCursor:         req.NextCursor,
	})
}

func (c *ControllerV1) SuggestKeywords(ctx context.Context, req *httpv1.SuggestKeywordsReq) (*httpv1.SuggestKeywordsRes, error) {
	return c.search.SuggestKeywords(ctx, &pb.SuggestKeywordsReq{
		Prefix: req.Prefix,
		Limit:  req.Limit,
	})
}

func (c *ControllerV1) BatchGetSpuCards(ctx context.Context, req *httpv1.BatchGetSpuCardsReq) (*httpv1.BatchGetSpuCardsRes, error) {
	return c.search.BatchGetSpuCards(ctx, &req.BatchGetSpuCardsReq)
}

func (c *ControllerV1) UpsertSpuDoc(ctx context.Context, req *httpv1.UpsertSpuDocReq) (*httpv1.UpsertSpuDocRes, error) {
	return c.search.UpsertSpuDoc(ctx, &req.UpsertSpuDocReq)
}

func (c *ControllerV1) PatchSpuStoreCategory(ctx context.Context, req *httpv1.PatchSpuStoreCategoryReq) (*httpv1.PatchSpuStoreCategoryRes, error) {
	return c.search.PatchSpuStoreCategory(ctx, &req.PatchSpuStoreCategoryReq)
}

func (c *ControllerV1) DeleteSpuDoc(ctx context.Context, req *httpv1.DeleteSpuDocReq) (*httpv1.DeleteSpuDocRes, error) {
	return c.search.DeleteSpuDoc(ctx, &req.DeleteSpuDocReq)
}

func (c *ControllerV1) RebuildIndex(ctx context.Context, req *httpv1.RebuildIndexReq) (*httpv1.RebuildIndexRes, error) {
	return c.search.RebuildIndex(ctx, &req.RebuildIndexReq)
}
