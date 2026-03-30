package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/search-svc/api/v1"
	"github.com/TsingpekTao/shopa/search-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// Controller 聚合 gRPC 接口层，只负责协议转发。
type Controller struct {
	v1.UnimplementedPublicSearchServiceServer
	v1.UnimplementedInternalSearchServiceServer
}

// Register 注册搜索服务的 gRPC 实现。
func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterPublicSearchServiceServer(s.Server, ctrl)
	v1.RegisterInternalSearchServiceServer(s.Server, ctrl)
}

// SearchProducts 调用 service 层执行分页检索。
func (*Controller) SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (res *v1.SearchProductsRes, err error) {
	return service.Search().SearchProducts(ctx, req)
}

// SuggestKeywords 调用 service 层执行联想词查询。
func (*Controller) SuggestKeywords(ctx context.Context, req *v1.SuggestKeywordsReq) (res *v1.SuggestKeywordsRes, err error) {
	return service.Search().SuggestKeywords(ctx, req)
}

// BatchGetSpuCards 调用 service 层执行批量卡片读取。
func (*Controller) BatchGetSpuCards(ctx context.Context, req *v1.BatchGetSpuCardsReq) (res *v1.BatchGetSpuCardsRes, err error) {
	return service.Search().BatchGetSpuCards(ctx, req)
}

// UpsertSpuDoc 调用 service 层执行文档写入。
func (*Controller) UpsertSpuDoc(ctx context.Context, req *v1.UpsertSpuDocReq) (res *v1.UpsertSpuDocRes, err error) {
	return service.Search().UpsertSpuDoc(ctx, req)
}

// DeleteSpuDoc 调用 service 层执行文档删除。
func (*Controller) DeleteSpuDoc(ctx context.Context, req *v1.DeleteSpuDocReq) (res *v1.DeleteSpuDocRes, err error) {
	return service.Search().DeleteSpuDoc(ctx, req)
}

// RebuildIndex 调用 service 层创建重建任务。
func (*Controller) RebuildIndex(ctx context.Context, req *v1.RebuildIndexReq) (res *v1.RebuildIndexRes, err error) {
	return service.Search().RebuildIndex(ctx, req)
}
