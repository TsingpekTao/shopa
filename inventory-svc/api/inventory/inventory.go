package inventory

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/inventory-svc/api/inventory/v1"
)

// IInventoryV1 defines HTTP handlers for inventory business APIs.
type IInventoryV1 interface {
	BatchAdjustMySkuStock(ctx context.Context, req *v1.BatchAdjustMySkuStockReq) (res *v1.BatchAdjustMySkuStockRes, err error)
	ReserveStock(ctx context.Context, req *v1.ReserveStockReq) (res *v1.ReserveStockRes, err error)
	ConfirmReservation(ctx context.Context, req *v1.ConfirmReservationReq) (res *v1.ConfirmReservationRes, err error)
	CancelReservation(ctx context.Context, req *v1.CancelReservationReq) (res *v1.CancelReservationRes, err error)
	BatchAdjustStockByAdmin(ctx context.Context, req *v1.BatchAdjustStockByAdminReq) (res *v1.BatchAdjustStockByAdminRes, err error)
	SetHotSku(ctx context.Context, req *v1.SetHotSkuReq) (res *v1.SetHotSkuRes, err error)
	GetSkuInventory(ctx context.Context, req *v1.GetSkuInventoryReq) (res *v1.GetSkuInventoryRes, err error)
	BatchGetSkuInventory(ctx context.Context, req *v1.BatchGetSkuInventoryReq) (res *v1.BatchGetSkuInventoryRes, err error)
	UpsertSkuContext(ctx context.Context, req *v1.UpsertSkuContextReq) (res *v1.UpsertSkuContextRes, err error)
}
