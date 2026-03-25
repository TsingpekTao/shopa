// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	ICatalog interface {
		// CreateProductDraft 创建商品草稿（SPU 主档），供商家在正式提审前先保存基础信息与图文资料。
		CreateProductDraft(ctx context.Context, req *v1.CreateProductDraftReq) (*v1.CreateProductDraftRes, error)
		// UpdateProductDraft 按 FieldMask 局部更新草稿字段，避免全量覆盖导致并发写冲突。
		UpdateProductDraft(ctx context.Context, req *v1.UpdateProductDraftReq) (*v1.UpdateProductDraftRes, error)
		// UpsertSkuDrafts 批量新增/更新 SKU 草稿明细，支持 replace_all 模式一次性替换整组 SKU。
		UpsertSkuDrafts(ctx context.Context, req *v1.UpsertSkuDraftsReq) (*v1.UpsertSkuDraftsRes, error)
		// SubmitProductReview 提交商品进入审核流程，仅允许可提审状态（如 DRAFT/OFF_SHELF）发起。
		SubmitProductReview(ctx context.Context, req *v1.SubmitProductReviewReq) (*v1.SubmitProductReviewRes, error)
		// ResubmitProductReview 针对驳回商品二次提审，通常用于补齐资质或修正违规信息后重提。
		ResubmitProductReview(ctx context.Context, req *v1.ResubmitProductReviewReq) (*v1.ResubmitProductReviewRes, error)
		// SetProductOnShelf 执行商品上架（支持立即上架或预约上架），审核通过后由商家主动触发。
		SetProductOnShelf(ctx context.Context, req *v1.SetProductOnShelfReq) (*v1.SetProductOnShelfRes, error)
		// SetProductOffShelf 执行商家主动下架，下架后前台不再展示但商品数据仍保留。
		SetProductOffShelf(ctx context.Context, req *v1.SetProductOffShelfReq) (*v1.SetProductOffShelfRes, error)
		// DeleteProductDraft 删除商品草稿（软删除），用于放弃未发布或已驳回的发品方案。
		DeleteProductDraft(ctx context.Context, req *v1.DeleteProductDraftReq) (*emptypb.Empty, error)
		// GetMyProduct 查询商家视角的单商品详情，包含审核信息与完整 SKU 结构。
		GetMyProduct(ctx context.Context, req *v1.GetMyProductReq) (*v1.GetMyProductRes, error)
		// ListMyProducts 分页查询商家商品列表，支持状态与关键词筛选。
		ListMyProducts(ctx context.Context, req *v1.ListMyProductsReq) (*v1.ListMyProductsRes, error)
		// ListReviewTasks 运营/审核后台分页查看审核任务队列与处理状态。
		ListReviewTasks(ctx context.Context, req *v1.ListReviewTasksReq) (*v1.ListReviewTasksRes, error)
		// GetReviewDetail 查看单个商品的审核详情与驳回原因，便于复核与追溯。
		GetReviewDetail(ctx context.Context, req *v1.GetReviewDetailReq) (*v1.GetReviewDetailRes, error)
		// ApproveProduct 审核通过商品，状态进入 APPROVED（待上架），不直接对买家可见。
		ApproveProduct(ctx context.Context, req *v1.ApproveProductReq) (*v1.ApproveProductRes, error)
		// RejectProduct 审核驳回商品，需记录驳回码与驳回说明，供商家后续整改。
		RejectProduct(ctx context.Context, req *v1.RejectProductReq) (*v1.RejectProductRes, error)
		// FreezeProduct 风控冻结商品，通常用于违规品处置，冻结后禁止继续售卖。
		FreezeProduct(ctx context.Context, req *v1.FreezeProductReq) (*v1.FreezeProductRes, error)
		// UnfreezeProduct 解冻商品，恢复到可继续运营的状态（由运营或风控确认后执行）。
		UnfreezeProduct(ctx context.Context, req *v1.UnfreezeProductReq) (*v1.UnfreezeProductRes, error)
		// ForceOffShelf 平台强制下架商品，用于紧急场景快速止损。
		ForceOffShelf(ctx context.Context, req *v1.ForceOffShelfReq) (*v1.ForceOffShelfRes, error)
		// GetProductDetail 买家侧查询商品详情，仅返回可售范围内的数据（如上架且可见）。
		GetProductDetail(ctx context.Context, req *v1.GetProductDetailReq) (*v1.GetProductDetailRes, error)
		// ListProducts 买家侧分页查询商品列表，承载类目页/频道页等读流量入口。
		ListProducts(ctx context.Context, req *v1.ListProductsReq) (*v1.ListProductsRes, error)
		// SearchProducts 买家侧关键词搜索接口，支持排序与筛选组合。
		SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error)
		// BatchGetSpuByNo 内部接口：按 spu_no 批量查询 SPU，供网关聚合或下游服务调用。
		BatchGetSpuByNo(ctx context.Context, req *v1.BatchGetSpuByNoReq) (*v1.BatchGetSpuByNoRes, error)
		// BatchGetSkuByNo 内部接口：按 sku_no 批量查询 SKU，适合订单预校验或购物车批量拉取。
		BatchGetSkuByNo(ctx context.Context, req *v1.BatchGetSkuByNoReq) (*v1.BatchGetSkuByNoRes, error)
		// GetSkuSnapshotForOrder 内部接口：为订单创建提供不可变 SKU 快照，防止后续改价穿透历史订单。
		GetSkuSnapshotForOrder(ctx context.Context, req *v1.GetSkuSnapshotForOrderReq) (*v1.GetSkuSnapshotForOrderRes, error)
		// UpsertSkuStockProjection 内部接口：接收库存事件并更新 SKU 库存投影，内置版本比较防止乱序覆盖。
		UpsertSkuStockProjection(ctx context.Context, req *v1.UpsertSkuStockProjectionReq) (*v1.UpsertSkuStockProjectionRes, error)
		// RecomputeSpuAggregation 内部接口：重算 SPU 聚合字段（价格区间/有货状态），用于列表页高性能读取。
		RecomputeSpuAggregation(ctx context.Context, req *v1.RecomputeSpuAggregationReq) (*v1.RecomputeSpuAggregationRes, error)
	}
)

var (
	localCatalog ICatalog
)

func Catalog() ICatalog {
	if localCatalog == nil {
		panic("implement not found for interface ICatalog, forgot register?")
	}
	return localCatalog
}

// RegisterCatalog 注册 Catalog 领域服务实现，供 controller 层通过 service.Catalog() 统一访问。
func RegisterCatalog(i ICatalog) {
	localCatalog = i
}
