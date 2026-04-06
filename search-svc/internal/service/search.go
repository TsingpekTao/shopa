// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/search-svc/api/v1"
)

type (
	ISearch interface {
		SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error)
		SuggestKeywords(ctx context.Context, req *v1.SuggestKeywordsReq) (*v1.SuggestKeywordsRes, error)
		BatchGetSpuCards(ctx context.Context, req *v1.BatchGetSpuCardsReq) (*v1.BatchGetSpuCardsRes, error)
		UpsertSpuDoc(ctx context.Context, req *v1.UpsertSpuDocReq) (*v1.UpsertSpuDocRes, error)
		PatchSpuStoreCategory(ctx context.Context, req *v1.PatchSpuStoreCategoryReq) (*v1.PatchSpuStoreCategoryRes, error)
		DeleteSpuDoc(ctx context.Context, req *v1.DeleteSpuDocReq) (*v1.DeleteSpuDocRes, error)
		RebuildIndex(ctx context.Context, req *v1.RebuildIndexReq) (*v1.RebuildIndexRes, error)
	}
)

var (
	localSearch ISearch
)

func Search() ISearch {
	if localSearch == nil {
		panic("implement not found for interface ISearch, forgot register?")
	}
	return localSearch
}

func RegisterSearch(i ISearch) {
	localSearch = i
}
