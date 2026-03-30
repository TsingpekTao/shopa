package cart

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/cart-svc/api/cart/v1"
)

// ICartV1 defines HTTP handlers for cart APIs.
type ICartV1 interface {
	AddItem(ctx context.Context, req *v1.AddItemReq) (res *v1.AddItemRes, err error)
	UpdateItemQty(ctx context.Context, req *v1.UpdateItemQtyReq) (res *v1.UpdateItemQtyRes, err error)
	ToggleItemChecked(ctx context.Context, req *v1.ToggleItemCheckedReq) (res *v1.ToggleItemCheckedRes, err error)
	BatchToggleItems(ctx context.Context, req *v1.BatchToggleItemsReq) (res *v1.BatchToggleItemsRes, err error)
	RemoveItems(ctx context.Context, req *v1.RemoveItemsReq) (res *v1.RemoveItemsRes, err error)
	ClearInvalidItems(ctx context.Context, req *v1.ClearInvalidItemsReq) (res *v1.ClearInvalidItemsRes, err error)
	GetMyCart(ctx context.Context, req *v1.GetMyCartReq) (res *v1.GetMyCartRes, err error)
	PrepareCheckout(ctx context.Context, req *v1.PrepareCheckoutReq) (res *v1.PrepareCheckoutRes, err error)
	ConsumeCheckoutToken(ctx context.Context, req *v1.ConsumeCheckoutTokenReq) (res *v1.ConsumeCheckoutTokenRes, err error)
	MarkItemsOrdered(ctx context.Context, req *v1.MarkItemsOrderedReq) (res *v1.MarkItemsOrderedRes, err error)
	BatchUpsertSkuProjection(ctx context.Context, req *v1.BatchUpsertSkuProjectionReq) (res *v1.BatchUpsertSkuProjectionRes, err error)
}
