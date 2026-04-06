package catalog

import (
	"context"
	"encoding/json"

	httpv1 "github.com/TsingpekTao/shopa/catalog-svc/api/catalog/v1"
	pb "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *ControllerV1) CreateProductDraft(ctx context.Context, req *httpv1.CreateProductDraftReq) (*httpv1.CreateProductDraftRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.CreateProductDraft(ctx, &pb.CreateProductDraftReq{
		ShopNo:              req.ShopNo,
		Title:               req.Title,
		SubTitle:            req.SubTitle,
		CategoryId:          req.CategoryId,
		BrandNo:             req.BrandNo,
		MainImageAssetIds:   req.MainImageAssetIds,
		DetailImageAssetIds: req.DetailImageAssetIds,
		SpuAttrs:            req.SpuAttrs,
		PublishTime:         req.PublishTime,
	})
}

func (c *ControllerV1) UpdateProductDraft(ctx context.Context, req *httpv1.UpdateProductDraftReq) (*httpv1.UpdateProductDraftRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.UpdateProductDraft(ctx, &pb.UpdateProductDraftReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		Patch:           req.Patch,
		UpdateMask:      req.UpdateMask,
	})
}

func (c *ControllerV1) UpsertSkuDrafts(ctx context.Context, req *httpv1.UpsertSkuDraftsReq) (*httpv1.UpsertSkuDraftsRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.UpsertSkuDrafts(ctx, &pb.UpsertSkuDraftsReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		Items:           req.Items,
		ReplaceAll:      req.ReplaceAll,
	})
}

func (c *ControllerV1) SubmitProductReview(ctx context.Context, req *httpv1.SubmitProductReviewReq) (*httpv1.SubmitProductReviewRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.SubmitProductReview(ctx, &pb.SubmitProductReviewReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		SubmitNote:      req.SubmitNote,
	})
}

func (c *ControllerV1) ResubmitProductReview(ctx context.Context, req *httpv1.ResubmitProductReviewReq) (*httpv1.ResubmitProductReviewRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.ResubmitProductReview(ctx, &pb.ResubmitProductReviewReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		SubmitNote:      req.SubmitNote,
	})
}

func (c *ControllerV1) SetProductOnShelf(ctx context.Context, req *httpv1.SetProductOnShelfReq) (*httpv1.SetProductOnShelfRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.SetProductOnShelf(ctx, &pb.SetProductOnShelfReq{
		SpuNo:                req.SpuNo,
		ExpectedVersion:      req.ExpectedVersion,
		EffectivePublishTime: req.EffectivePublishTime,
	})
}

func (c *ControllerV1) SetProductOffShelf(ctx context.Context, req *httpv1.SetProductOffShelfReq) (*httpv1.SetProductOffShelfRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.SetProductOffShelf(ctx, &pb.SetProductOffShelfReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		ReasonCode:      req.ReasonCode,
	})
}

func (c *ControllerV1) DeleteProductDraft(ctx context.Context, req *httpv1.DeleteProductDraftReq) (*emptypb.Empty, error) {
	hydrateJSONBody(ctx, req)
	_, err := c.catalog.DeleteProductDraft(ctx, &pb.DeleteProductDraftReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		ReasonCode:      req.ReasonCode,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c *ControllerV1) GetMyProduct(ctx context.Context, req *httpv1.GetMyProductReq) (*httpv1.GetMyProductRes, error) {
	return c.catalog.GetMyProduct(ctx, &pb.GetMyProductReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ListMyProducts(ctx context.Context, req *httpv1.ListMyProductsReq) (*httpv1.ListMyProductsRes, error) {
	return c.catalog.ListMyProducts(ctx, &pb.ListMyProductsReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Statuses: req.Statuses,
		Keyword:  req.Keyword,
	})
}

func (c *ControllerV1) ListReviewTasks(ctx context.Context, req *httpv1.ListReviewTasksReq) (*httpv1.ListReviewTasksRes, error) {
	return c.catalog.ListReviewTasks(ctx, &pb.ListReviewTasksReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Statuses: req.Statuses,
		Keyword:  req.Keyword,
	})
}

func (c *ControllerV1) GetReviewDetail(ctx context.Context, req *httpv1.GetReviewDetailReq) (*httpv1.GetReviewDetailRes, error) {
	return c.catalog.GetReviewDetail(ctx, &pb.GetReviewDetailReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ApproveProduct(ctx context.Context, req *httpv1.ApproveProductReq) (*httpv1.ApproveProductRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.ApproveProduct(ctx, &pb.ApproveProductReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		ReviewComment:   req.ReviewComment,
	})
}

func (c *ControllerV1) RejectProduct(ctx context.Context, req *httpv1.RejectProductReq) (*httpv1.RejectProductRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.RejectProduct(ctx, &pb.RejectProductReq{
		SpuNo:            req.SpuNo,
		ExpectedVersion:  req.ExpectedVersion,
		RejectReasonCode: req.RejectReasonCode,
		RejectComment:    req.RejectComment,
	})
}

func (c *ControllerV1) FreezeProduct(ctx context.Context, req *httpv1.FreezeProductReq) (*httpv1.FreezeProductRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.FreezeProduct(ctx, &pb.FreezeProductReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		ReasonCode:      req.ReasonCode,
		Reason:          req.Reason,
	})
}

func (c *ControllerV1) UnfreezeProduct(ctx context.Context, req *httpv1.UnfreezeProductReq) (*httpv1.UnfreezeProductRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.UnfreezeProduct(ctx, &pb.UnfreezeProductReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
	})
}

func (c *ControllerV1) ForceOffShelf(ctx context.Context, req *httpv1.ForceOffShelfReq) (*httpv1.ForceOffShelfRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.ForceOffShelf(ctx, &pb.ForceOffShelfReq{
		SpuNo:           req.SpuNo,
		ExpectedVersion: req.ExpectedVersion,
		ReasonCode:      req.ReasonCode,
		Reason:          req.Reason,
	})
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

func (c *ControllerV1) UpsertSkuStockProjectionInternal(ctx context.Context, req *httpv1.UpsertSkuStockProjectionInternalReq) (*httpv1.UpsertSkuStockProjectionInternalRes, error) {
	hydrateJSONBody(ctx, req)
	return c.catalog.UpsertSkuStockProjection(ctx, &pb.UpsertSkuStockProjectionReq{
		Items: req.Items,
	})
}

func hydrateJSONBody(ctx context.Context, dst interface{}) {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return
	}
	_ = gconv.Scan(req.GetMap(), dst)
	if j, err := req.GetJson(); err == nil && j != nil {
		_ = j.Scan(dst)
	}
	body := req.GetBody()
	if len(body) == 0 {
		return
	}
	_ = json.Unmarshal(body, dst)
}

var _ = emptypb.Empty{}
