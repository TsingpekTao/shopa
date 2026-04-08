package seller

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) ListRefundBatches(ctx context.Context, req *v1.ListRefundBatchesReq) (*v1.ListRefundBatchesRes, error) {
	return service.Bff().ListSellerRefundBatches(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), req)
}

func (c *ControllerV1) GetRefundBatchDetail(ctx context.Context, req *v1.GetRefundBatchDetailReq) (*v1.GetRefundBatchDetailRes, error) {
	return service.Bff().GetSellerRefundBatchDetail(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), req)
}

func (c *ControllerV1) ApproveRefundBatch(ctx context.Context, req *v1.ApproveRefundBatchReq) (*v1.ApproveRefundBatchRes, error) {
	return service.Bff().ApproveSellerRefundBatch(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), &req.ApproveRefundBatchReq)
}

func (c *ControllerV1) RejectRefundBatch(ctx context.Context, req *v1.RejectRefundBatchReq) (*v1.RejectRefundBatchRes, error) {
	return service.Bff().RejectSellerRefundBatch(ctx, extractAccessTokenFromRequest(g.RequestFromCtx(ctx)), &req.RejectRefundBatchReq)
}
