// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 你可以删除这些注释，并按需手动维护该接口文件。
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	// ICatalog 定义商品领域服务接口，覆盖商家草稿、审核、上下架、买家查询与库存投影能力。
	ICatalog interface {
		// CreateProductDraft 创建商品草稿（SPU）。
		CreateProductDraft(ctx context.Context, req *v1.CreateProductDraftReq) (*v1.CreateProductDraftRes, error)
		// UpdateProductDraft 按 FieldMask 更新商品草稿字段。
		UpdateProductDraft(ctx context.Context, req *v1.UpdateProductDraftReq) (*v1.UpdateProductDraftRes, error)
		// UpsertSkuDrafts 批量新增或更新 SKU 草稿。
		UpsertSkuDrafts(ctx context.Context, req *v1.UpsertSkuDraftsReq) (*v1.UpsertSkuDraftsRes, error)
		// SubmitProductReview 提交商品审核。
		SubmitProductReview(ctx context.Context, req *v1.SubmitProductReviewReq) (*v1.SubmitProductReviewRes, error)
		// ResubmitProductReview 驳回后重新提交审核。
		ResubmitProductReview(ctx context.Context, req *v1.ResubmitProductReviewReq) (*v1.ResubmitProductReviewRes, error)
		// SetProductOnShelf 将商品上架。
		SetProductOnShelf(ctx context.Context, req *v1.SetProductOnShelfReq) (*v1.SetProductOnShelfRes, error)
		// SetProductOffShelf 将商品下架。
		SetProductOffShelf(ctx context.Context, req *v1.SetProductOffShelfReq) (*v1.SetProductOffShelfRes, error)
		// DeleteProductDraft 删除商品草稿。
		DeleteProductDraft(ctx context.Context, req *v1.DeleteProductDraftReq) (*emptypb.Empty, error)
		// GetMyProduct 查询商家单个商品详情。
		GetMyProduct(ctx context.Context, req *v1.GetMyProductReq) (*v1.GetMyProductRes, error)
		// ListMyProducts 分页查询商家商品列表。
		ListMyProducts(ctx context.Context, req *v1.ListMyProductsReq) (*v1.ListMyProductsRes, error)
		// ListReviewTasks 分页查询审核任务列表。
		ListReviewTasks(ctx context.Context, req *v1.ListReviewTasksReq) (*v1.ListReviewTasksRes, error)
		// GetReviewDetail 查询审核任务详情。
		GetReviewDetail(ctx context.Context, req *v1.GetReviewDetailReq) (*v1.GetReviewDetailRes, error)
		// ApproveProduct 通过商品审核。
		ApproveProduct(ctx context.Context, req *v1.ApproveProductReq) (*v1.ApproveProductRes, error)
		// RejectProduct 驳回商品审核。
		RejectProduct(ctx context.Context, req *v1.RejectProductReq) (*v1.RejectProductRes, error)
		// FreezeProduct 冻结商品。
		FreezeProduct(ctx context.Context, req *v1.FreezeProductReq) (*v1.FreezeProductRes, error)
		// UnfreezeProduct 解冻商品。
		UnfreezeProduct(ctx context.Context, req *v1.UnfreezeProductReq) (*v1.UnfreezeProductRes, error)
		// ForceOffShelf 平台强制下架商品。
		ForceOffShelf(ctx context.Context, req *v1.ForceOffShelfReq) (*v1.ForceOffShelfRes, error)
		// GetProductDetail 查询买家侧商品详情。
		GetProductDetail(ctx context.Context, req *v1.GetProductDetailReq) (*v1.GetProductDetailRes, error)
		// ListProducts 分页查询商品列表。
		ListProducts(ctx context.Context, req *v1.ListProductsReq) (*v1.ListProductsRes, error)
		// SearchProducts 按关键词搜索商品。
		SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error)
		// BatchGetSpuByNo 批量查询 SPU。
		BatchGetSpuByNo(ctx context.Context, req *v1.BatchGetSpuByNoReq) (*v1.BatchGetSpuByNoRes, error)
		// BatchGetSkuByNo 批量查询 SKU。
		BatchGetSkuByNo(ctx context.Context, req *v1.BatchGetSkuByNoReq) (*v1.BatchGetSkuByNoRes, error)
		// GetSkuSnapshotForOrder 获取下单所需 SKU 快照。
		GetSkuSnapshotForOrder(ctx context.Context, req *v1.GetSkuSnapshotForOrderReq) (*v1.GetSkuSnapshotForOrderRes, error)
		// UpsertSkuStockProjection 更新 SKU 库存状态投影。
		UpsertSkuStockProjection(ctx context.Context, req *v1.UpsertSkuStockProjectionReq) (*v1.UpsertSkuStockProjectionRes, error)
		// RecomputeSpuAggregation 重算 SPU 聚合字段。
		RecomputeSpuAggregation(ctx context.Context, req *v1.RecomputeSpuAggregationReq) (*v1.RecomputeSpuAggregationRes, error)
	}
)

var (
	localCatalog ICatalog
)

// Catalog 返回已注册的商品服务实现。
func Catalog() ICatalog {
	if localCatalog == nil {
		panic("implement not found for interface ICatalog, forgot register?")
	}
	return localCatalog
}

// RegisterCatalog 注册商品服务实现，通常在 internal/logic 的 init 中调用。
func RegisterCatalog(i ICatalog) {
	localCatalog = i
}
