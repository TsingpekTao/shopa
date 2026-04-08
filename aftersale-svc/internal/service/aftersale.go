// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
)

type (
	IAfterSale interface {
		CreateAfterSale(ctx context.Context, req *v1.CreateAfterSaleReq) (*v1.CreateAfterSaleRes, error)
		CancelAfterSale(ctx context.Context, req *v1.CancelAfterSaleReq) (*v1.CancelAfterSaleRes, error)
		GetMyAfterSaleDetail(ctx context.Context, req *v1.GetMyAfterSaleDetailReq) (*v1.GetMyAfterSaleDetailRes, error)
		ListMyAfterSales(ctx context.Context, req *v1.ListMyAfterSalesReq) (*v1.ListMyAfterSalesRes, error)
		ApplyRefundBatch(ctx context.Context, req *v1.ApplyRefundBatchReq) (*v1.ApplyRefundBatchRes, error)
		ListMyRefundBatches(ctx context.Context, req *v1.ListMyRefundBatchesReq) (*v1.ListMyRefundBatchesRes, error)
		GetMyRefundBatchDetail(ctx context.Context, req *v1.GetMyRefundBatchDetailReq) (*v1.GetMyRefundBatchDetailRes, error)
		CancelRefundBatch(ctx context.Context, req *v1.CancelRefundBatchReq) (*v1.CancelRefundBatchRes, error)
		ListShopAfterSales(ctx context.Context, req *v1.ListShopAfterSalesReq) (*v1.ListShopAfterSalesRes, error)
		GetShopAfterSaleDetail(ctx context.Context, req *v1.GetShopAfterSaleDetailReq) (*v1.GetShopAfterSaleDetailRes, error)
		ApproveAfterSale(ctx context.Context, req *v1.ApproveAfterSaleReq) (*v1.ApproveAfterSaleRes, error)
		RejectAfterSale(ctx context.Context, req *v1.RejectAfterSaleReq) (*v1.RejectAfterSaleRes, error)
		ListShopRefundBatches(ctx context.Context, req *v1.ListShopRefundBatchesReq) (*v1.ListShopRefundBatchesRes, error)
		GetShopRefundBatchDetail(ctx context.Context, req *v1.GetShopRefundBatchDetailReq) (*v1.GetShopRefundBatchDetailRes, error)
		ApproveRefundBatch(ctx context.Context, req *v1.ApproveRefundBatchReq) (*v1.ApproveRefundBatchRes, error)
		RejectRefundBatch(ctx context.Context, req *v1.RejectRefundBatchReq) (*v1.RejectRefundBatchRes, error)
		ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (*v1.ExecuteRefundTaskRes, error)
		RetryRefundTask(ctx context.Context, req *v1.RetryRefundTaskReq) (*v1.RetryRefundTaskRes, error)
		GetAfterSaleSnapshotByNo(ctx context.Context, req *v1.GetAfterSaleSnapshotByNoReq) (*v1.GetAfterSaleSnapshotByNoRes, error)
	}
)

var (
	localAfterSale IAfterSale
)

func AfterSale() IAfterSale {
	if localAfterSale == nil {
		panic("implement not found for interface IAfterSale, forgot register?")
	}
	return localAfterSale
}

func RegisterAfterSale(i IAfterSale) {
	localAfterSale = i
}
