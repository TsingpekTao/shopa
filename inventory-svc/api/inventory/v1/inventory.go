package v1

import (
	pb "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type BatchAdjustMySkuStockReq struct {
	g.Meta `path:"/v1/inventory/seller/stock:adjust" method:"post" tags:"Inventory-Seller" summary:"Batch adjust my sku stock"`
	pb.BatchAdjustMySkuStockReq
}
type BatchAdjustMySkuStockRes = pb.BatchAdjustMySkuStockRes

type ReserveStockReq struct {
	g.Meta `path:"/v1/inventory/order/reservations" method:"post" tags:"Inventory-Order" summary:"Reserve stock"`
	pb.ReserveStockReq
}
type ReserveStockRes = pb.ReserveStockRes

type ConfirmReservationReq struct {
	g.Meta `path:"/v1/inventory/order/reservations:confirm" method:"post" tags:"Inventory-Order" summary:"Confirm reservation"`
	pb.ConfirmReservationReq
}
type ConfirmReservationRes = pb.ConfirmReservationRes

type CancelReservationReq struct {
	g.Meta `path:"/v1/inventory/order/reservations:cancel" method:"post" tags:"Inventory-Order" summary:"Cancel reservation"`
	pb.CancelReservationReq
}
type CancelReservationRes = pb.CancelReservationRes

type BatchAdjustStockByAdminReq struct {
	g.Meta `path:"/v1/inventory/admin/stock:adjust" method:"post" tags:"Inventory-Admin" summary:"Batch adjust stock by admin"`
	pb.BatchAdjustStockByAdminReq
}
type BatchAdjustStockByAdminRes = pb.BatchAdjustStockByAdminRes

type SetHotSkuReq struct {
	g.Meta `path:"/v1/inventory/admin/sku/hot:set" method:"post" tags:"Inventory-Admin" summary:"Set hot sku"`
	pb.SetHotSkuReq
}
type SetHotSkuRes = pb.SetHotSkuRes

type GetSkuInventoryReq struct {
	g.Meta `path:"/v1/inventory/sku/{sku_no}" method:"get" tags:"Inventory-Query" summary:"Get sku inventory"`
	SkuNo  string `json:"sku_no"`
}
type GetSkuInventoryRes = pb.GetSkuInventoryRes

type BatchGetSkuInventoryReq struct {
	g.Meta `path:"/v1/inventory/sku:batch-get" method:"post" tags:"Inventory-Query" summary:"Batch get sku inventory"`
	pb.BatchGetSkuInventoryReq
}
type BatchGetSkuInventoryRes = pb.BatchGetSkuInventoryRes

type UpsertSkuContextReq struct {
	g.Meta `path:"/v1/inventory/internal/sku-context:upsert" method:"post" tags:"Inventory-Internal" summary:"Upsert sku context"`
	pb.UpsertSkuContextReq
}
type UpsertSkuContextRes = pb.UpsertSkuContextRes
