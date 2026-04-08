package aftersale

import (
	"context"

	aftersalev1 "github.com/TsingpekTao/shopa/aftersale-svc/api/aftersale/v1"
	pb "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
)

func (c *ControllerV1) CreateAfterSale(ctx context.Context, req *aftersalev1.CreateAfterSaleReq) (*aftersalev1.CreateAfterSaleRes, error) {
	return c.afterSale.CreateAfterSale(ctx, &req.CreateAfterSaleReq)
}

func (c *ControllerV1) CreateAfterSaleAlias(ctx context.Context, req *aftersalev1.CreateAfterSaleAliasReq) (*aftersalev1.CreateAfterSaleAliasRes, error) {
	return c.afterSale.CreateAfterSale(ctx, &req.CreateAfterSaleReq)
}

func (c *ControllerV1) CancelAfterSale(ctx context.Context, req *aftersalev1.CancelAfterSaleReq) (*aftersalev1.CancelAfterSaleRes, error) {
	return c.afterSale.CancelAfterSale(ctx, &req.CancelAfterSaleReq)
}

func (c *ControllerV1) CancelAfterSaleAlias(ctx context.Context, req *aftersalev1.CancelAfterSaleAliasReq) (*aftersalev1.CancelAfterSaleAliasRes, error) {
	return c.afterSale.CancelAfterSale(ctx, &req.CancelAfterSaleReq)
}

func (c *ControllerV1) GetMyAfterSaleDetail(ctx context.Context, req *aftersalev1.GetMyAfterSaleDetailReq) (*aftersalev1.GetMyAfterSaleDetailRes, error) {
	return c.afterSale.GetMyAfterSaleDetail(ctx, &pb.GetMyAfterSaleDetailReq{AfterSaleNo: req.AfterSaleNo})
}

func (c *ControllerV1) GetMyAfterSaleDetailAlias(ctx context.Context, req *aftersalev1.GetMyAfterSaleDetailAliasReq) (*aftersalev1.GetMyAfterSaleDetailAliasRes, error) {
	return c.afterSale.GetMyAfterSaleDetail(ctx, &pb.GetMyAfterSaleDetailReq{AfterSaleNo: req.AfterSaleNo})
}

func (c *ControllerV1) ListMyAfterSales(ctx context.Context, req *aftersalev1.ListMyAfterSalesReq) (*aftersalev1.ListMyAfterSalesRes, error) {
	return c.afterSale.ListMyAfterSales(ctx, &req.ListMyAfterSalesReq)
}

func (c *ControllerV1) ListMyAfterSalesAlias(ctx context.Context, req *aftersalev1.ListMyAfterSalesAliasReq) (*aftersalev1.ListMyAfterSalesAliasRes, error) {
	return c.afterSale.ListMyAfterSales(ctx, &pb.ListMyAfterSalesReq{
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
		Statuses:   req.Statuses,
	})
}

func (c *ControllerV1) ApplyRefundBatch(ctx context.Context, req *aftersalev1.ApplyRefundBatchReq) (*aftersalev1.ApplyRefundBatchRes, error) {
	return c.afterSale.ApplyRefundBatch(ctx, &req.ApplyRefundBatchReq)
}

func (c *ControllerV1) ListMyRefundBatches(ctx context.Context, req *aftersalev1.ListMyRefundBatchesReq) (*aftersalev1.ListMyRefundBatchesRes, error) {
	return c.afterSale.ListMyRefundBatches(ctx, &pb.ListMyRefundBatchesReq{
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
		Statuses:   req.Statuses,
	})
}

func (c *ControllerV1) GetMyRefundBatchDetail(ctx context.Context, req *aftersalev1.GetMyRefundBatchDetailReq) (*aftersalev1.GetMyRefundBatchDetailRes, error) {
	return c.afterSale.GetMyRefundBatchDetail(ctx, &pb.GetMyRefundBatchDetailReq{RefundBatchNo: req.RefundBatchNo})
}

func (c *ControllerV1) CancelRefundBatch(ctx context.Context, req *aftersalev1.CancelRefundBatchReq) (*aftersalev1.CancelRefundBatchRes, error) {
	return c.afterSale.CancelRefundBatch(ctx, &req.CancelRefundBatchReq)
}

func (c *ControllerV1) ListShopAfterSales(ctx context.Context, req *aftersalev1.ListShopAfterSalesReq) (*aftersalev1.ListShopAfterSalesRes, error) {
	return c.afterSale.ListShopAfterSales(ctx, &req.ListShopAfterSalesReq)
}

func (c *ControllerV1) ListShopAfterSalesAlias(ctx context.Context, req *aftersalev1.ListShopAfterSalesAliasReq) (*aftersalev1.ListShopAfterSalesAliasRes, error) {
	return c.afterSale.ListShopAfterSales(ctx, &pb.ListShopAfterSalesReq{
		ShopNo:     req.ShopNo,
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
		Statuses:   req.Statuses,
	})
}

func (c *ControllerV1) GetShopAfterSaleDetail(ctx context.Context, req *aftersalev1.GetShopAfterSaleDetailReq) (*aftersalev1.GetShopAfterSaleDetailRes, error) {
	return c.afterSale.GetShopAfterSaleDetail(ctx, &pb.GetShopAfterSaleDetailReq{AfterSaleNo: req.AfterSaleNo})
}

func (c *ControllerV1) GetShopAfterSaleDetailAlias(ctx context.Context, req *aftersalev1.GetShopAfterSaleDetailAliasReq) (*aftersalev1.GetShopAfterSaleDetailAliasRes, error) {
	return c.afterSale.GetShopAfterSaleDetail(ctx, &pb.GetShopAfterSaleDetailReq{AfterSaleNo: req.AfterSaleNo})
}

func (c *ControllerV1) ApproveAfterSale(ctx context.Context, req *aftersalev1.ApproveAfterSaleReq) (*aftersalev1.ApproveAfterSaleRes, error) {
	return c.afterSale.ApproveAfterSale(ctx, &req.ApproveAfterSaleReq)
}

func (c *ControllerV1) ApproveAfterSaleAlias(ctx context.Context, req *aftersalev1.ApproveAfterSaleAliasReq) (*aftersalev1.ApproveAfterSaleAliasRes, error) {
	return c.afterSale.ApproveAfterSale(ctx, &req.ApproveAfterSaleReq)
}

func (c *ControllerV1) RejectAfterSale(ctx context.Context, req *aftersalev1.RejectAfterSaleReq) (*aftersalev1.RejectAfterSaleRes, error) {
	return c.afterSale.RejectAfterSale(ctx, &req.RejectAfterSaleReq)
}

func (c *ControllerV1) RejectAfterSaleAlias(ctx context.Context, req *aftersalev1.RejectAfterSaleAliasReq) (*aftersalev1.RejectAfterSaleAliasRes, error) {
	return c.afterSale.RejectAfterSale(ctx, &req.RejectAfterSaleReq)
}

func (c *ControllerV1) ListShopRefundBatches(ctx context.Context, req *aftersalev1.ListShopRefundBatchesReq) (*aftersalev1.ListShopRefundBatchesRes, error) {
	return c.afterSale.ListShopRefundBatches(ctx, &pb.ListShopRefundBatchesReq{
		ShopNo:     req.ShopNo,
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
		Statuses:   req.Statuses,
	})
}

func (c *ControllerV1) GetShopRefundBatchDetail(ctx context.Context, req *aftersalev1.GetShopRefundBatchDetailReq) (*aftersalev1.GetShopRefundBatchDetailRes, error) {
	return c.afterSale.GetShopRefundBatchDetail(ctx, &pb.GetShopRefundBatchDetailReq{
		RefundBatchNo: req.RefundBatchNo,
		ShopNo:        req.ShopNo,
	})
}

func (c *ControllerV1) ApproveRefundBatch(ctx context.Context, req *aftersalev1.ApproveRefundBatchReq) (*aftersalev1.ApproveRefundBatchRes, error) {
	return c.afterSale.ApproveRefundBatch(ctx, &req.ApproveRefundBatchReq)
}

func (c *ControllerV1) RejectRefundBatch(ctx context.Context, req *aftersalev1.RejectRefundBatchReq) (*aftersalev1.RejectRefundBatchRes, error) {
	return c.afterSale.RejectRefundBatch(ctx, &req.RejectRefundBatchReq)
}

func (c *ControllerV1) ExecuteRefundTask(ctx context.Context, req *aftersalev1.ExecuteRefundTaskReq) (*aftersalev1.ExecuteRefundTaskRes, error) {
	return c.afterSale.ExecuteRefundTask(ctx, &req.ExecuteRefundTaskReq)
}

func (c *ControllerV1) RetryRefundTask(ctx context.Context, req *aftersalev1.RetryRefundTaskReq) (*aftersalev1.RetryRefundTaskRes, error) {
	return c.afterSale.RetryRefundTask(ctx, &req.RetryRefundTaskReq)
}

func (c *ControllerV1) GetAfterSaleSnapshotByNo(ctx context.Context, req *aftersalev1.GetAfterSaleSnapshotByNoReq) (*aftersalev1.GetAfterSaleSnapshotByNoRes, error) {
	return c.afterSale.GetAfterSaleSnapshotByNo(ctx, &pb.GetAfterSaleSnapshotByNoReq{AfterSaleNo: req.AfterSaleNo})
}
