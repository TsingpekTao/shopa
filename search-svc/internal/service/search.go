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
		// SearchProducts 根据关键字与过滤条件返回商品卡片列表。
		SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error)
		// SuggestKeywords 根据前缀返回建议词。
		SuggestKeywords(ctx context.Context, req *v1.SuggestKeywordsReq) (*v1.SuggestKeywordsRes, error)
		// BatchGetSpuCards 批量读取 spu 卡片信息。
		BatchGetSpuCards(ctx context.Context, req *v1.BatchGetSpuCardsReq) (*v1.BatchGetSpuCardsRes, error)
		// UpsertSpuDoc 写入或更新搜索文档。
		UpsertSpuDoc(ctx context.Context, req *v1.UpsertSpuDocReq) (*v1.UpsertSpuDocRes, error)
		// DeleteSpuDoc 删除指定 spu 的搜索文档。
		DeleteSpuDoc(ctx context.Context, req *v1.DeleteSpuDocReq) (*v1.DeleteSpuDocRes, error)
		// RebuildIndex 触发索引重建任务。
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
