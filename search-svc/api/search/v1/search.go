package v1

import (
	pb "github.com/TsingpekTao/shopa/search-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type SearchProductsReq struct {
	g.Meta             `path:"/v1/search/products" method:"get" tags:"Search" summary:"Search products"`
	Query              string `json:"query" in:"query"`
	ShopNo             string `json:"shopNo" in:"query"`
	CategoryNo         string `json:"categoryNo" in:"query"`
	StoreCategoryId    uint64 `json:"storeCategoryId" in:"query"`
	StoreCategoryLevel int32  `json:"storeCategoryLevel" in:"query"`
	SortCode           string `json:"sortCode" in:"query"`
	PageSize           int32  `json:"pageSize" in:"query"`
	NextCursor         string `json:"nextCursor" in:"query"`
}

type SearchProductsRes = pb.SearchProductsRes

type SuggestKeywordsReq struct {
	g.Meta `path:"/v1/search/suggest" method:"get" tags:"Search" summary:"Suggest keywords"`
	Prefix string `json:"prefix" in:"query"`
	Limit  int32  `json:"limit" in:"query"`
}

type SuggestKeywordsRes = pb.SuggestKeywordsRes

type BatchGetSpuCardsReq struct {
	g.Meta `path:"/v1/search/spu-cards/batch" method:"post" tags:"Search" summary:"Batch get spu cards"`
	pb.BatchGetSpuCardsReq
}

type BatchGetSpuCardsRes = pb.BatchGetSpuCardsRes

type UpsertSpuDocReq struct {
	g.Meta `path:"/v1/internal/search/docs/upsert" method:"post" tags:"SearchInternal" summary:"Upsert search spu doc"`
	pb.UpsertSpuDocReq
}

type UpsertSpuDocRes = pb.UpsertSpuDocRes

type PatchSpuStoreCategoryReq struct {
	g.Meta `path:"/v1/internal/search/docs/store-category" method:"post" tags:"SearchInternal" summary:"Patch search spu store category"`
	pb.PatchSpuStoreCategoryReq
}

type PatchSpuStoreCategoryRes = pb.PatchSpuStoreCategoryRes

type DeleteSpuDocReq struct {
	g.Meta `path:"/v1/internal/search/docs/delete" method:"post" tags:"SearchInternal" summary:"Delete search spu doc"`
	pb.DeleteSpuDocReq
}

type DeleteSpuDocRes = pb.DeleteSpuDocRes

type RebuildIndexReq struct {
	g.Meta `path:"/v1/internal/search/rebuild" method:"post" tags:"SearchInternal" summary:"Trigger search index rebuild"`
	pb.RebuildIndexReq
}

type RebuildIndexRes = pb.RebuildIndexRes
