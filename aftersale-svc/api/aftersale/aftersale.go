package aftersale

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/aftersale-svc/api/aftersale/v1"
)

type IAftersaleV1 interface {
	CreateAfterSale(ctx context.Context, req *v1.CreateAfterSaleReq) (res *v1.CreateAfterSaleRes, err error)
	CreateAfterSaleAlias(ctx context.Context, req *v1.CreateAfterSaleAliasReq) (res *v1.CreateAfterSaleAliasRes, err error)
	CancelAfterSale(ctx context.Context, req *v1.CancelAfterSaleReq) (res *v1.CancelAfterSaleRes, err error)
	CancelAfterSaleAlias(ctx context.Context, req *v1.CancelAfterSaleAliasReq) (res *v1.CancelAfterSaleAliasRes, err error)
	GetMyAfterSaleDetail(ctx context.Context, req *v1.GetMyAfterSaleDetailReq) (res *v1.GetMyAfterSaleDetailRes, err error)
	GetMyAfterSaleDetailAlias(ctx context.Context, req *v1.GetMyAfterSaleDetailAliasReq) (res *v1.GetMyAfterSaleDetailAliasRes, err error)
	ListMyAfterSales(ctx context.Context, req *v1.ListMyAfterSalesReq) (res *v1.ListMyAfterSalesRes, err error)
	ListMyAfterSalesAlias(ctx context.Context, req *v1.ListMyAfterSalesAliasReq) (res *v1.ListMyAfterSalesAliasRes, err error)
	ApplyRefundBatch(ctx context.Context, req *v1.ApplyRefundBatchReq) (res *v1.ApplyRefundBatchRes, err error)
	ListMyRefundBatches(ctx context.Context, req *v1.ListMyRefundBatchesReq) (res *v1.ListMyRefundBatchesRes, err error)
	GetMyRefundBatchDetail(ctx context.Context, req *v1.GetMyRefundBatchDetailReq) (res *v1.GetMyRefundBatchDetailRes, err error)
	CancelRefundBatch(ctx context.Context, req *v1.CancelRefundBatchReq) (res *v1.CancelRefundBatchRes, err error)
	ListShopAfterSales(ctx context.Context, req *v1.ListShopAfterSalesReq) (res *v1.ListShopAfterSalesRes, err error)
	ListShopAfterSalesAlias(ctx context.Context, req *v1.ListShopAfterSalesAliasReq) (res *v1.ListShopAfterSalesAliasRes, err error)
	GetShopAfterSaleDetail(ctx context.Context, req *v1.GetShopAfterSaleDetailReq) (res *v1.GetShopAfterSaleDetailRes, err error)
	GetShopAfterSaleDetailAlias(ctx context.Context, req *v1.GetShopAfterSaleDetailAliasReq) (res *v1.GetShopAfterSaleDetailAliasRes, err error)
	ApproveAfterSale(ctx context.Context, req *v1.ApproveAfterSaleReq) (res *v1.ApproveAfterSaleRes, err error)
	ApproveAfterSaleAlias(ctx context.Context, req *v1.ApproveAfterSaleAliasReq) (res *v1.ApproveAfterSaleAliasRes, err error)
	RejectAfterSale(ctx context.Context, req *v1.RejectAfterSaleReq) (res *v1.RejectAfterSaleRes, err error)
	RejectAfterSaleAlias(ctx context.Context, req *v1.RejectAfterSaleAliasReq) (res *v1.RejectAfterSaleAliasRes, err error)
	ListShopRefundBatches(ctx context.Context, req *v1.ListShopRefundBatchesReq) (res *v1.ListShopRefundBatchesRes, err error)
	GetShopRefundBatchDetail(ctx context.Context, req *v1.GetShopRefundBatchDetailReq) (res *v1.GetShopRefundBatchDetailRes, err error)
	ApproveRefundBatch(ctx context.Context, req *v1.ApproveRefundBatchReq) (res *v1.ApproveRefundBatchRes, err error)
	RejectRefundBatch(ctx context.Context, req *v1.RejectRefundBatchReq) (res *v1.RejectRefundBatchRes, err error)
	ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (res *v1.ExecuteRefundTaskRes, err error)
	RetryRefundTask(ctx context.Context, req *v1.RetryRefundTaskReq) (res *v1.RetryRefundTaskRes, err error)
	GetAfterSaleSnapshotByNo(ctx context.Context, req *v1.GetAfterSaleSnapshotByNoReq) (res *v1.GetAfterSaleSnapshotByNoRes, err error)
}
