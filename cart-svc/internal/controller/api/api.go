package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/cart-svc/api/v1"
	"github.com/TsingpekTao/shopa/cart-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedBuyerCartServiceServer
	v1.UnimplementedInternalCartServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterBuyerCartServiceServer(s.Server, ctrl)
	v1.RegisterInternalCartServiceServer(s.Server, ctrl)
}

func (*Controller) AddItem(ctx context.Context, req *v1.AddItemReq) (*v1.AddItemRes, error) {
	return service.Cart().AddItem(ctx, req)
}

func (*Controller) UpdateItemQty(ctx context.Context, req *v1.UpdateItemQtyReq) (*v1.UpdateItemQtyRes, error) {
	return service.Cart().UpdateItemQty(ctx, req)
}

func (*Controller) ToggleItemChecked(ctx context.Context, req *v1.ToggleItemCheckedReq) (*v1.ToggleItemCheckedRes, error) {
	return service.Cart().ToggleItemChecked(ctx, req)
}

func (*Controller) BatchToggleItems(ctx context.Context, req *v1.BatchToggleItemsReq) (*v1.BatchToggleItemsRes, error) {
	return service.Cart().BatchToggleItems(ctx, req)
}

func (*Controller) RemoveItems(ctx context.Context, req *v1.RemoveItemsReq) (*v1.RemoveItemsRes, error) {
	return service.Cart().RemoveItems(ctx, req)
}

func (*Controller) ClearInvalidItems(ctx context.Context, req *v1.ClearInvalidItemsReq) (*v1.ClearInvalidItemsRes, error) {
	return service.Cart().ClearInvalidItems(ctx, req)
}

func (*Controller) GetMyCart(ctx context.Context, req *v1.GetMyCartReq) (*v1.GetMyCartRes, error) {
	return service.Cart().GetMyCart(ctx, req)
}

func (*Controller) PrepareCheckout(ctx context.Context, req *v1.PrepareCheckoutReq) (*v1.PrepareCheckoutRes, error) {
	return service.Cart().PrepareCheckout(ctx, req)
}

func (*Controller) ConsumeCheckoutToken(ctx context.Context, req *v1.ConsumeCheckoutTokenReq) (*v1.ConsumeCheckoutTokenRes, error) {
	return service.Cart().ConsumeCheckoutToken(ctx, req)
}

func (*Controller) MarkItemsOrdered(ctx context.Context, req *v1.MarkItemsOrderedReq) (*v1.MarkItemsOrderedRes, error) {
	return service.Cart().MarkItemsOrdered(ctx, req)
}

func (*Controller) BatchUpsertSkuProjection(ctx context.Context, req *v1.BatchUpsertSkuProjectionReq) (*v1.BatchUpsertSkuProjectionRes, error) {
	return service.Cart().BatchUpsertSkuProjection(ctx, req)
}
