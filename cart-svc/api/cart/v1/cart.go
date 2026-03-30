package v1

import (
	pb "github.com/TsingpekTao/shopa/cart-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type AddItemReq struct {
	g.Meta `path:"/v1/cart/items:add" method:"post" tags:"Cart-Buyer" summary:"Add item to cart"`
	pb.AddItemReq
}
type AddItemRes = pb.AddItemRes

type UpdateItemQtyReq struct {
	g.Meta `path:"/v1/cart/items:qty" method:"post" tags:"Cart-Buyer" summary:"Update cart item qty"`
	pb.UpdateItemQtyReq
}
type UpdateItemQtyRes = pb.UpdateItemQtyRes

type ToggleItemCheckedReq struct {
	g.Meta `path:"/v1/cart/items:check" method:"post" tags:"Cart-Buyer" summary:"Toggle cart item checked status"`
	pb.ToggleItemCheckedReq
}
type ToggleItemCheckedRes = pb.ToggleItemCheckedRes

type BatchToggleItemsReq struct {
	g.Meta `path:"/v1/cart/items:batch-check" method:"post" tags:"Cart-Buyer" summary:"Batch toggle cart item checked status"`
	pb.BatchToggleItemsReq
}
type BatchToggleItemsRes = pb.BatchToggleItemsRes

type RemoveItemsReq struct {
	g.Meta `path:"/v1/cart/items:remove" method:"post" tags:"Cart-Buyer" summary:"Remove cart items"`
	pb.RemoveItemsReq
}
type RemoveItemsRes = pb.RemoveItemsRes

type ClearInvalidItemsReq struct {
	g.Meta `path:"/v1/cart/items:clear-invalid" method:"post" tags:"Cart-Buyer" summary:"Clear invalid cart items"`
	pb.ClearInvalidItemsReq
}
type ClearInvalidItemsRes = pb.ClearInvalidItemsRes

type GetMyCartReq struct {
	g.Meta      `path:"/v1/cart/me" method:"get" tags:"Cart-Buyer" summary:"Get my cart"`
	OnlyChecked bool `json:"only_checked"`
}
type GetMyCartRes = pb.GetMyCartRes

type PrepareCheckoutReq struct {
	g.Meta `path:"/v1/cart/checkout:prepare" method:"post" tags:"Cart-Buyer" summary:"Prepare checkout snapshot token"`
	pb.PrepareCheckoutReq
}
type PrepareCheckoutRes = pb.PrepareCheckoutRes

type ConsumeCheckoutTokenReq struct {
	g.Meta `path:"/v1/cart/internal/checkout:consume" method:"post" tags:"Cart-Internal" summary:"Consume checkout token atomically"`
	pb.ConsumeCheckoutTokenReq
}
type ConsumeCheckoutTokenRes = pb.ConsumeCheckoutTokenRes

type MarkItemsOrderedReq struct {
	g.Meta `path:"/v1/cart/internal/items:ordered" method:"post" tags:"Cart-Internal" summary:"Mark cart items as ordered"`
	pb.MarkItemsOrderedReq
}
type MarkItemsOrderedRes = pb.MarkItemsOrderedRes

type BatchUpsertSkuProjectionReq struct {
	g.Meta `path:"/v1/cart/internal/sku-projection:batch-upsert" method:"post" tags:"Cart-Internal" summary:"Batch upsert sku projection"`
	pb.BatchUpsertSkuProjectionReq
}
type BatchUpsertSkuProjectionRes = pb.BatchUpsertSkuProjectionRes
