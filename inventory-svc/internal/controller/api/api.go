package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	"github.com/TsingpekTao/shopa/inventory-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedSellerInventoryServiceServer
	v1.UnimplementedOrderInventoryServiceServer
	v1.UnimplementedAdminInventoryServiceServer
	v1.UnimplementedInternalInventoryServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterSellerInventoryServiceServer(s.Server, ctrl)
	v1.RegisterOrderInventoryServiceServer(s.Server, ctrl)
	v1.RegisterAdminInventoryServiceServer(s.Server, ctrl)
	v1.RegisterInternalInventoryServiceServer(s.Server, ctrl)
}

func (*Controller) BatchAdjustMySkuStock(ctx context.Context, req *v1.BatchAdjustMySkuStockReq) (*v1.BatchAdjustMySkuStockRes, error) {
	return service.Inventory().BatchAdjustMySkuStock(ctx, req)
}

func (*Controller) ReserveStock(ctx context.Context, req *v1.ReserveStockReq) (*v1.ReserveStockRes, error) {
	return service.Inventory().ReserveStock(ctx, req)
}

func (*Controller) ConfirmReservation(ctx context.Context, req *v1.ConfirmReservationReq) (*v1.ConfirmReservationRes, error) {
	return service.Inventory().ConfirmReservation(ctx, req)
}

func (*Controller) CancelReservation(ctx context.Context, req *v1.CancelReservationReq) (*v1.CancelReservationRes, error) {
	return service.Inventory().CancelReservation(ctx, req)
}

func (*Controller) BatchAdjustStockByAdmin(ctx context.Context, req *v1.BatchAdjustStockByAdminReq) (*v1.BatchAdjustStockByAdminRes, error) {
	return service.Inventory().BatchAdjustStockByAdmin(ctx, req)
}

func (*Controller) SetHotSku(ctx context.Context, req *v1.SetHotSkuReq) (*v1.SetHotSkuRes, error) {
	return service.Inventory().SetHotSku(ctx, req)
}

func (*Controller) GetSkuInventory(ctx context.Context, req *v1.GetSkuInventoryReq) (*v1.GetSkuInventoryRes, error) {
	return service.Inventory().GetSkuInventory(ctx, req)
}

func (*Controller) BatchGetSkuInventory(ctx context.Context, req *v1.BatchGetSkuInventoryReq) (*v1.BatchGetSkuInventoryRes, error) {
	return service.Inventory().BatchGetSkuInventory(ctx, req)
}

func (*Controller) UpsertSkuContext(ctx context.Context, req *v1.UpsertSkuContextReq) (*v1.UpsertSkuContextRes, error) {
	return service.Inventory().UpsertSkuContext(ctx, req)
}
