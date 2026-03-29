package inventory

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	pb "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
)

func (c *ControllerV1) BatchAdjustMySkuStock(ctx context.Context, req *v1.BatchAdjustMySkuStockReq) (*v1.BatchAdjustMySkuStockRes, error) {
	return c.inventory.BatchAdjustMySkuStock(ctx, req)
}

func (c *ControllerV1) ReserveStock(ctx context.Context, req *v1.ReserveStockReq) (*v1.ReserveStockRes, error) {
	return c.inventory.ReserveStock(ctx, req)
}

func (c *ControllerV1) ConfirmReservation(ctx context.Context, req *v1.ConfirmReservationReq) (*v1.ConfirmReservationRes, error) {
	return c.inventory.ConfirmReservation(ctx, req)
}

func (c *ControllerV1) CancelReservation(ctx context.Context, req *v1.CancelReservationReq) (*v1.CancelReservationRes, error) {
	return c.inventory.CancelReservation(ctx, req)
}

func (c *ControllerV1) BatchAdjustStockByAdmin(ctx context.Context, req *v1.BatchAdjustStockByAdminReq) (*v1.BatchAdjustStockByAdminRes, error) {
	return c.inventory.BatchAdjustStockByAdmin(ctx, req)
}

func (c *ControllerV1) SetHotSku(ctx context.Context, req *v1.SetHotSkuReq) (*v1.SetHotSkuRes, error) {
	return c.inventory.SetHotSku(ctx, req)
}

func (c *ControllerV1) GetSkuInventory(ctx context.Context, req *v1.GetSkuInventoryReq) (*v1.GetSkuInventoryRes, error) {
	return c.inventory.GetSkuInventory(ctx, &pb.GetSkuInventoryReq{SkuNo: req.SkuNo})
}

func (c *ControllerV1) BatchGetSkuInventory(ctx context.Context, req *v1.BatchGetSkuInventoryReq) (*v1.BatchGetSkuInventoryRes, error) {
	return c.inventory.BatchGetSkuInventory(ctx, req)
}

func (c *ControllerV1) UpsertSkuContext(ctx context.Context, req *v1.UpsertSkuContextReq) (*v1.UpsertSkuContextRes, error) {
	return c.inventory.UpsertSkuContext(ctx, req)
}

