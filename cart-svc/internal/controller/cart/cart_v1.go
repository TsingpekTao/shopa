package cart

import (
	"context"
	"encoding/json"

	cartv1 "github.com/TsingpekTao/shopa/cart-svc/api/cart/v1"
	pbv1 "github.com/TsingpekTao/shopa/cart-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) AddItem(ctx context.Context, req *cartv1.AddItemReq) (*cartv1.AddItemRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.AddItem(ctx, &pbv1.AddItemReq{
		SkuNo:          req.SkuNo,
		SpuNo:          req.SpuNo,
		ShopNo:         req.ShopNo,
		Qty:            req.Qty,
		Checked:        req.Checked,
		IdempotencyKey: req.IdempotencyKey,
	})
}

func (c *ControllerV1) UpdateItemQty(ctx context.Context, req *cartv1.UpdateItemQtyReq) (*cartv1.UpdateItemQtyRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.UpdateItemQty(ctx, &pbv1.UpdateItemQtyReq{
		SkuNo:          req.SkuNo,
		Qty:            req.Qty,
		IdempotencyKey: req.IdempotencyKey,
	})
}

func (c *ControllerV1) ToggleItemChecked(ctx context.Context, req *cartv1.ToggleItemCheckedReq) (*cartv1.ToggleItemCheckedRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.ToggleItemChecked(ctx, &pbv1.ToggleItemCheckedReq{
		SkuNo:   req.SkuNo,
		Checked: req.Checked,
	})
}

func (c *ControllerV1) BatchToggleItems(ctx context.Context, req *cartv1.BatchToggleItemsReq) (*cartv1.BatchToggleItemsRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.BatchToggleItems(ctx, &pbv1.BatchToggleItemsReq{
		SkuNos:  req.SkuNos,
		Checked: req.Checked,
	})
}

func (c *ControllerV1) RemoveItems(ctx context.Context, req *cartv1.RemoveItemsReq) (*cartv1.RemoveItemsRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.RemoveItems(ctx, &pbv1.RemoveItemsReq{
		SkuNos: req.SkuNos,
	})
}

func (c *ControllerV1) ClearInvalidItems(ctx context.Context, req *cartv1.ClearInvalidItemsReq) (*cartv1.ClearInvalidItemsRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.ClearInvalidItems(ctx, &pbv1.ClearInvalidItemsReq{})
}

func (c *ControllerV1) GetMyCart(ctx context.Context, req *cartv1.GetMyCartReq) (*cartv1.GetMyCartRes, error) {
	return c.cart.GetMyCart(ctx, &pbv1.GetMyCartReq{OnlyChecked: req.OnlyChecked})
}

func (c *ControllerV1) PrepareCheckout(ctx context.Context, req *cartv1.PrepareCheckoutReq) (*cartv1.PrepareCheckoutRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.PrepareCheckout(ctx, &pbv1.PrepareCheckoutReq{
		Scope:     req.Scope,
		SkuNos:    req.SkuNos,
		AddressId: req.AddressId,
	})
}

func (c *ControllerV1) ConsumeCheckoutToken(ctx context.Context, req *cartv1.ConsumeCheckoutTokenReq) (*cartv1.ConsumeCheckoutTokenRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.ConsumeCheckoutToken(ctx, &pbv1.ConsumeCheckoutTokenReq{
		CheckoutToken: req.CheckoutToken,
		UserId:        req.UserId,
	})
}

func (c *ControllerV1) MarkItemsOrdered(ctx context.Context, req *cartv1.MarkItemsOrderedReq) (*cartv1.MarkItemsOrderedRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.MarkItemsOrdered(ctx, &pbv1.MarkItemsOrderedReq{
		UserId:  req.UserId,
		OrderNo: req.OrderNo,
		SkuNos:  req.SkuNos,
	})
}

func (c *ControllerV1) BatchUpsertSkuProjection(ctx context.Context, req *cartv1.BatchUpsertSkuProjectionReq) (*cartv1.BatchUpsertSkuProjectionRes, error) {
	hydrateJSONBody(ctx, req)
	return c.cart.BatchUpsertSkuProjection(ctx, &pbv1.BatchUpsertSkuProjectionReq{
		Items: req.Items,
	})
}

func hydrateJSONBody(ctx context.Context, dst interface{}) {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return
	}
	if err := req.Parse(dst); err == nil {
		return
	}
	_ = gconv.Scan(req.GetMap(), dst)
	if j, err := req.GetJson(); err == nil && j != nil {
		_ = j.Scan(dst)
	}
	body := req.GetBody()
	if len(body) == 0 {
		return
	}
	_ = json.Unmarshal(body, dst)
}
