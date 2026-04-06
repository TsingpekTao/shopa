package search

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/search-svc/api/search/v1"
)

type ISearchV1 interface {
	SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (res *v1.SearchProductsRes, err error)
	SuggestKeywords(ctx context.Context, req *v1.SuggestKeywordsReq) (res *v1.SuggestKeywordsRes, err error)
	BatchGetSpuCards(ctx context.Context, req *v1.BatchGetSpuCardsReq) (res *v1.BatchGetSpuCardsRes, err error)
	UpsertSpuDoc(ctx context.Context, req *v1.UpsertSpuDocReq) (res *v1.UpsertSpuDocRes, err error)
	PatchSpuStoreCategory(ctx context.Context, req *v1.PatchSpuStoreCategoryReq) (res *v1.PatchSpuStoreCategoryRes, err error)
	DeleteSpuDoc(ctx context.Context, req *v1.DeleteSpuDocReq) (res *v1.DeleteSpuDocRes, err error)
	RebuildIndex(ctx context.Context, req *v1.RebuildIndexReq) (res *v1.RebuildIndexRes, err error)
}
