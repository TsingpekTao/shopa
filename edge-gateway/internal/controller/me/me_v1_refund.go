package me

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) PreviewRefund(ctx context.Context, req *v1.PreviewRefundReq) (*v1.PreviewRefundRes, error) {
	return service.Bff().BuildBuyerRefundPreview(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), req)
}

func (c *ControllerV1) ApplyRefundBatch(ctx context.Context, req *v1.ApplyRefundBatchReq) (*v1.ApplyRefundBatchRes, error) {
	return service.Bff().ApplyBuyerRefundBatch(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), &req.ApplyRefundBatchReq)
}

func (c *ControllerV1) ListRefundBatches(ctx context.Context, req *v1.ListRefundBatchesReq) (*v1.ListRefundBatchesRes, error) {
	return service.Bff().ListBuyerRefundBatches(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), req)
}

func (c *ControllerV1) GetRefundBatchDetail(ctx context.Context, req *v1.GetRefundBatchDetailReq) (*v1.GetRefundBatchDetailRes, error) {
	return service.Bff().GetBuyerRefundBatchDetail(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), req)
}

func (c *ControllerV1) CancelRefundBatch(ctx context.Context, req *v1.CancelRefundBatchReq) (*v1.CancelRefundBatchRes, error) {
	return service.Bff().CancelBuyerRefundBatch(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), &req.CancelRefundBatchReq)
}
