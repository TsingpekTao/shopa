package inventory

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/inventory-svc/api/inventory/v1"
)

// IInventoryV1 瀹氫箟搴撳瓨涓氬姟 API 鐨?HTTP 澶勭悊鍣ㄣ€?
type IInventoryV1 interface {
	BatchAdjustMySkuStock(ctx context.Context, req *httpv1.BatchAdjustMySkuStockReq) (res *httpv1.BatchAdjustMySkuStockRes, err error)
	ReserveStock(ctx context.Context, req *httpv1.ReserveStockReq) (res *httpv1.ReserveStockRes, err error)
	ConfirmReservation(ctx context.Context, req *httpv1.ConfirmReservationReq) (res *httpv1.ConfirmReservationRes, err error)
	CancelReservation(ctx context.Context, req *httpv1.CancelReservationReq) (res *httpv1.CancelReservationRes, err error)
	BatchAdjustStockByAdmin(ctx context.Context, req *httpv1.BatchAdjustStockByAdminReq) (res *httpv1.BatchAdjustStockByAdminRes, err error)
	SetHotSku(ctx context.Context, req *httpv1.SetHotSkuReq) (res *httpv1.SetHotSkuRes, err error)
	GetSkuInventory(ctx context.Context, req *httpv1.GetSkuInventoryReq) (res *httpv1.GetSkuInventoryRes, err error)
	BatchGetSkuInventory(ctx context.Context, req *httpv1.BatchGetSkuInventoryReq) (res *httpv1.BatchGetSkuInventoryRes, err error)
	UpsertSkuContext(ctx context.Context, req *httpv1.UpsertSkuContextReq) (res *httpv1.UpsertSkuContextRes, err error)
}
