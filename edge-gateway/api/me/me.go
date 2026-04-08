package me

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/me/v1"
)

type IMeV1 interface {
	GetOverview(ctx context.Context, req *v1.GetOverviewReq) (res *v1.GetOverviewRes, err error)
	PreviewRefund(ctx context.Context, req *v1.PreviewRefundReq) (res *v1.PreviewRefundRes, err error)
	ApplyRefundBatch(ctx context.Context, req *v1.ApplyRefundBatchReq) (res *v1.ApplyRefundBatchRes, err error)
	ListRefundBatches(ctx context.Context, req *v1.ListRefundBatchesReq) (res *v1.ListRefundBatchesRes, err error)
	GetRefundBatchDetail(ctx context.Context, req *v1.GetRefundBatchDetailReq) (res *v1.GetRefundBatchDetailRes, err error)
	CancelRefundBatch(ctx context.Context, req *v1.CancelRefundBatchReq) (res *v1.CancelRefundBatchRes, err error)
}
