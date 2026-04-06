package inventory

import (
	"context"
	"encoding/json"

	httpv1 "github.com/TsingpekTao/shopa/inventory-svc/api/inventory/v1"
	pb "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) BatchAdjustMySkuStock(ctx context.Context, req *httpv1.BatchAdjustMySkuStockReq) (*httpv1.BatchAdjustMySkuStockRes, error) {
	hydrateInventoryJSONBody(ctx, &req.BatchAdjustMySkuStockReq)
	return c.inventory.BatchAdjustMySkuStock(ctx, &req.BatchAdjustMySkuStockReq)
}

func (c *ControllerV1) ReserveStock(ctx context.Context, req *httpv1.ReserveStockReq) (*httpv1.ReserveStockRes, error) {
	hydrateInventoryJSONBody(ctx, &req.ReserveStockReq)
	return c.inventory.ReserveStock(ctx, &req.ReserveStockReq)
}

func (c *ControllerV1) ConfirmReservation(ctx context.Context, req *httpv1.ConfirmReservationReq) (*httpv1.ConfirmReservationRes, error) {
	hydrateInventoryJSONBody(ctx, &req.ConfirmReservationReq)
	return c.inventory.ConfirmReservation(ctx, &req.ConfirmReservationReq)
}

func (c *ControllerV1) CancelReservation(ctx context.Context, req *httpv1.CancelReservationReq) (*httpv1.CancelReservationRes, error) {
	hydrateInventoryJSONBody(ctx, &req.CancelReservationReq)
	return c.inventory.CancelReservation(ctx, &req.CancelReservationReq)
}

func (c *ControllerV1) BatchAdjustStockByAdmin(ctx context.Context, req *httpv1.BatchAdjustStockByAdminReq) (*httpv1.BatchAdjustStockByAdminRes, error) {
	hydrateInventoryJSONBody(ctx, &req.BatchAdjustStockByAdminReq)
	return c.inventory.BatchAdjustStockByAdmin(ctx, &req.BatchAdjustStockByAdminReq)
}

func (c *ControllerV1) SetHotSku(ctx context.Context, req *httpv1.SetHotSkuReq) (*httpv1.SetHotSkuRes, error) {
	hydrateInventoryJSONBody(ctx, &req.SetHotSkuReq)
	return c.inventory.SetHotSku(ctx, &req.SetHotSkuReq)
}

func (c *ControllerV1) GetSkuInventory(ctx context.Context, req *httpv1.GetSkuInventoryReq) (*httpv1.GetSkuInventoryRes, error) {
	return c.inventory.GetSkuInventory(ctx, &pb.GetSkuInventoryReq{SkuNo: req.SkuNo})
}

func (c *ControllerV1) BatchGetSkuInventory(ctx context.Context, req *httpv1.BatchGetSkuInventoryReq) (*httpv1.BatchGetSkuInventoryRes, error) {
	hydrateInventoryJSONBody(ctx, &req.BatchGetSkuInventoryReq)
	return c.inventory.BatchGetSkuInventory(ctx, &req.BatchGetSkuInventoryReq)
}

func (c *ControllerV1) UpsertSkuContext(ctx context.Context, req *httpv1.UpsertSkuContextReq) (*httpv1.UpsertSkuContextRes, error) {
	hydrateInventoryJSONBody(ctx, &req.UpsertSkuContextReq)
	return c.inventory.UpsertSkuContext(ctx, &req.UpsertSkuContextReq)
}

func hydrateInventoryJSONBody(ctx context.Context, dst interface{}) {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return
	}
	body := req.GetBody()
	if len(body) == 0 {
		return
	}
	_ = json.Unmarshal(body, dst)
}
