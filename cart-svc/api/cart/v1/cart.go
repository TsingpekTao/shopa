package v1

import (
	pb "github.com/TsingpekTao/shopa/cart-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type AddItemReq struct {
	g.Meta         `path:"/v1/cart/items:add" method:"post" tags:"Cart-Buyer" summary:"Add item to cart"`
	SkuNo          string `json:"sku_no" p:"sku_no"`
	SpuNo          string `json:"spu_no" p:"spu_no"`
	ShopNo         string `json:"shop_no" p:"shop_no"`
	Qty            uint32 `json:"qty" p:"qty"`
	Checked        bool   `json:"checked" p:"checked"`
	IdempotencyKey string `json:"idempotency_key" p:"idempotency_key"`
}
type AddItemRes = pb.AddItemRes

type UpdateItemQtyReq struct {
	g.Meta         `path:"/v1/cart/items:qty" method:"post" tags:"Cart-Buyer" summary:"Update cart item qty"`
	SkuNo          string `json:"sku_no" p:"sku_no"`
	Qty            uint32 `json:"qty" p:"qty"`
	IdempotencyKey string `json:"idempotency_key" p:"idempotency_key"`
}
type UpdateItemQtyRes = pb.UpdateItemQtyRes

type ToggleItemCheckedReq struct {
	g.Meta  `path:"/v1/cart/items:check" method:"post" tags:"Cart-Buyer" summary:"Toggle cart item checked status"`
	SkuNo   string `json:"sku_no" p:"sku_no"`
	Checked bool   `json:"checked" p:"checked"`
}
type ToggleItemCheckedRes = pb.ToggleItemCheckedRes

type BatchToggleItemsReq struct {
	g.Meta  `path:"/v1/cart/items:batch-check" method:"post" tags:"Cart-Buyer" summary:"Batch toggle cart item checked status"`
	SkuNos  []string `json:"sku_nos" p:"sku_nos"`
	Checked bool     `json:"checked" p:"checked"`
}
type BatchToggleItemsRes = pb.BatchToggleItemsRes

type RemoveItemsReq struct {
	g.Meta `path:"/v1/cart/items:remove" method:"post" tags:"Cart-Buyer" summary:"Remove cart items"`
	SkuNos []string `json:"sku_nos" p:"sku_nos"`
}
type RemoveItemsRes = pb.RemoveItemsRes

type ClearInvalidItemsReq struct {
	g.Meta `path:"/v1/cart/items:clear-invalid" method:"post" tags:"Cart-Buyer" summary:"Clear invalid cart items"`
}
type ClearInvalidItemsRes = pb.ClearInvalidItemsRes

type GetMyCartReq struct {
	g.Meta      `path:"/v1/cart/me" method:"get" tags:"Cart-Buyer" summary:"Get my cart"`
	OnlyChecked bool `json:"only_checked" p:"only_checked"`
}
type GetMyCartRes = pb.GetMyCartRes

type PrepareCheckoutReq struct {
	g.Meta    `path:"/v1/cart/checkout:prepare" method:"post" tags:"Cart-Buyer" summary:"Prepare checkout snapshot token"`
	Scope     pb.CheckoutScope `json:"scope" p:"scope"`
	SkuNos    []string         `json:"sku_nos" p:"sku_nos"`
	AddressId uint64           `json:"address_id" p:"address_id"`
}
type PrepareCheckoutRes = pb.PrepareCheckoutRes

type ConsumeCheckoutTokenReq struct {
	g.Meta        `path:"/v1/cart/internal/checkout:consume" method:"post" tags:"Cart-Internal" summary:"Consume checkout token atomically"`
	CheckoutToken string `json:"checkout_token" p:"checkout_token"`
	UserId        uint64 `json:"user_id" p:"user_id"`
}
type ConsumeCheckoutTokenRes = pb.ConsumeCheckoutTokenRes

type MarkItemsOrderedReq struct {
	g.Meta  `path:"/v1/cart/internal/items:ordered" method:"post" tags:"Cart-Internal" summary:"Mark cart items as ordered"`
	UserId  uint64   `json:"user_id" p:"user_id"`
	OrderNo string   `json:"order_no" p:"order_no"`
	SkuNos  []string `json:"sku_nos" p:"sku_nos"`
}
type MarkItemsOrderedRes = pb.MarkItemsOrderedRes

type BatchUpsertSkuProjectionReq struct {
	g.Meta `path:"/v1/cart/internal/sku-projection:batch-upsert" method:"post" tags:"Cart-Internal" summary:"Batch upsert sku projection"`
	Items  []*pb.SkuProjectionItem `json:"items" p:"items"`
}
type BatchUpsertSkuProjectionRes = pb.BatchUpsertSkuProjectionRes
