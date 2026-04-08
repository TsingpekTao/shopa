package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedBuyerAfterSaleServiceServer
	v1.UnimplementedSellerAfterSaleServiceServer
	v1.UnimplementedInternalAfterSaleServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterBuyerAfterSaleServiceServer(s.Server, ctrl)
	v1.RegisterSellerAfterSaleServiceServer(s.Server, ctrl)
	v1.RegisterInternalAfterSaleServiceServer(s.Server, ctrl)
}

func (*Controller) CreateAfterSale(ctx context.Context, req *v1.CreateAfterSaleReq) (res *v1.CreateAfterSaleRes, err error) {
	return service.AfterSale().CreateAfterSale(ctx, req)
}

func (*Controller) CancelAfterSale(ctx context.Context, req *v1.CancelAfterSaleReq) (res *v1.CancelAfterSaleRes, err error) {
	return service.AfterSale().CancelAfterSale(ctx, req)
}

func (*Controller) GetMyAfterSaleDetail(ctx context.Context, req *v1.GetMyAfterSaleDetailReq) (res *v1.GetMyAfterSaleDetailRes, err error) {
	return service.AfterSale().GetMyAfterSaleDetail(ctx, req)
}

func (*Controller) ListMyAfterSales(ctx context.Context, req *v1.ListMyAfterSalesReq) (res *v1.ListMyAfterSalesRes, err error) {
	return service.AfterSale().ListMyAfterSales(ctx, req)
}

func (*Controller) ApplyRefundBatch(ctx context.Context, req *v1.ApplyRefundBatchReq) (res *v1.ApplyRefundBatchRes, err error) {
	return service.AfterSale().ApplyRefundBatch(ctx, req)
}

func (*Controller) ListMyRefundBatches(ctx context.Context, req *v1.ListMyRefundBatchesReq) (res *v1.ListMyRefundBatchesRes, err error) {
	return service.AfterSale().ListMyRefundBatches(ctx, req)
}

func (*Controller) GetMyRefundBatchDetail(ctx context.Context, req *v1.GetMyRefundBatchDetailReq) (res *v1.GetMyRefundBatchDetailRes, err error) {
	return service.AfterSale().GetMyRefundBatchDetail(ctx, req)
}

func (*Controller) CancelRefundBatch(ctx context.Context, req *v1.CancelRefundBatchReq) (res *v1.CancelRefundBatchRes, err error) {
	return service.AfterSale().CancelRefundBatch(ctx, req)
}

func (*Controller) ListShopAfterSales(ctx context.Context, req *v1.ListShopAfterSalesReq) (res *v1.ListShopAfterSalesRes, err error) {
	return service.AfterSale().ListShopAfterSales(ctx, req)
}

func (*Controller) GetShopAfterSaleDetail(ctx context.Context, req *v1.GetShopAfterSaleDetailReq) (res *v1.GetShopAfterSaleDetailRes, err error) {
	return service.AfterSale().GetShopAfterSaleDetail(ctx, req)
}

func (*Controller) ApproveAfterSale(ctx context.Context, req *v1.ApproveAfterSaleReq) (res *v1.ApproveAfterSaleRes, err error) {
	return service.AfterSale().ApproveAfterSale(ctx, req)
}

func (*Controller) RejectAfterSale(ctx context.Context, req *v1.RejectAfterSaleReq) (res *v1.RejectAfterSaleRes, err error) {
	return service.AfterSale().RejectAfterSale(ctx, req)
}

func (*Controller) ListShopRefundBatches(ctx context.Context, req *v1.ListShopRefundBatchesReq) (res *v1.ListShopRefundBatchesRes, err error) {
	return service.AfterSale().ListShopRefundBatches(ctx, req)
}

func (*Controller) GetShopRefundBatchDetail(ctx context.Context, req *v1.GetShopRefundBatchDetailReq) (res *v1.GetShopRefundBatchDetailRes, err error) {
	return service.AfterSale().GetShopRefundBatchDetail(ctx, req)
}

func (*Controller) ApproveRefundBatch(ctx context.Context, req *v1.ApproveRefundBatchReq) (res *v1.ApproveRefundBatchRes, err error) {
	return service.AfterSale().ApproveRefundBatch(ctx, req)
}

func (*Controller) RejectRefundBatch(ctx context.Context, req *v1.RejectRefundBatchReq) (res *v1.RejectRefundBatchRes, err error) {
	return service.AfterSale().RejectRefundBatch(ctx, req)
}

func (*Controller) ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (res *v1.ExecuteRefundTaskRes, err error) {
	return service.AfterSale().ExecuteRefundTask(ctx, req)
}

func (*Controller) RetryRefundTask(ctx context.Context, req *v1.RetryRefundTaskReq) (res *v1.RetryRefundTaskRes, err error) {
	return service.AfterSale().RetryRefundTask(ctx, req)
}

func (*Controller) GetAfterSaleSnapshotByNo(ctx context.Context, req *v1.GetAfterSaleSnapshotByNoReq) (res *v1.GetAfterSaleSnapshotByNoRes, err error) {
	return service.AfterSale().GetAfterSaleSnapshotByNo(ctx, req)
}
