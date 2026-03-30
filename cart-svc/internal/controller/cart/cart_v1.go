package cart

import (
	"context"

	cartv1 "github.com/TsingpekTao/shopa/cart-svc/api/cart/v1"
	pbv1 "github.com/TsingpekTao/shopa/cart-svc/api/v1"
)

func (c *ControllerV1) AddItem(ctx context.Context, req *cartv1.AddItemReq) (*cartv1.AddItemRes, error) {
	return c.cart.AddItem(ctx, &req.AddItemReq)
}

func (c *ControllerV1) UpdateItemQty(ctx context.Context, req *cartv1.UpdateItemQtyReq) (*cartv1.UpdateItemQtyRes, error) {
	return c.cart.UpdateItemQty(ctx, &req.UpdateItemQtyReq)
}

func (c *ControllerV1) ToggleItemChecked(ctx context.Context, req *cartv1.ToggleItemCheckedReq) (*cartv1.ToggleItemCheckedRes, error) {
	return c.cart.ToggleItemChecked(ctx, &req.ToggleItemCheckedReq)
}

func (c *ControllerV1) BatchToggleItems(ctx context.Context, req *cartv1.BatchToggleItemsReq) (*cartv1.BatchToggleItemsRes, error) {
	return c.cart.BatchToggleItems(ctx, &req.BatchToggleItemsReq)
}

func (c *ControllerV1) RemoveItems(ctx context.Context, req *cartv1.RemoveItemsReq) (*cartv1.RemoveItemsRes, error) {
	return c.cart.RemoveItems(ctx, &req.RemoveItemsReq)
}

func (c *ControllerV1) ClearInvalidItems(ctx context.Context, req *cartv1.ClearInvalidItemsReq) (*cartv1.ClearInvalidItemsRes, error) {
	return c.cart.ClearInvalidItems(ctx, &pbv1.ClearInvalidItemsReq{})
}

func (c *ControllerV1) GetMyCart(ctx context.Context, req *cartv1.GetMyCartReq) (*cartv1.GetMyCartRes, error) {
	return c.cart.GetMyCart(ctx, &pbv1.GetMyCartReq{OnlyChecked: req.OnlyChecked})
}

func (c *ControllerV1) PrepareCheckout(ctx context.Context, req *cartv1.PrepareCheckoutReq) (*cartv1.PrepareCheckoutRes, error) {
	return c.cart.PrepareCheckout(ctx, &req.PrepareCheckoutReq)
}

func (c *ControllerV1) ConsumeCheckoutToken(ctx context.Context, req *cartv1.ConsumeCheckoutTokenReq) (*cartv1.ConsumeCheckoutTokenRes, error) {
	return c.cart.ConsumeCheckoutToken(ctx, &req.ConsumeCheckoutTokenReq)
}

func (c *ControllerV1) MarkItemsOrdered(ctx context.Context, req *cartv1.MarkItemsOrderedReq) (*cartv1.MarkItemsOrderedRes, error) {
	return c.cart.MarkItemsOrdered(ctx, &req.MarkItemsOrderedReq)
}

func (c *ControllerV1) BatchUpsertSkuProjection(ctx context.Context, req *cartv1.BatchUpsertSkuProjectionReq) (*cartv1.BatchUpsertSkuProjectionRes, error) {
	return c.cart.BatchUpsertSkuProjection(ctx, &req.BatchUpsertSkuProjectionReq)
}
