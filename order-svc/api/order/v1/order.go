package v1

import (
	pb "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreateOrderFromCartReq struct {
	g.Meta `path:"/v1/order/buyer/from-cart" method:"post" tags:"Order-Buyer" summary:"Create order from checkout token"`
	pb.CreateOrderFromCartReq
}

type CreateOrderFromCartRes = pb.CreateOrderFromCartRes

type CreateOrderBuyNowReq struct {
	g.Meta `path:"/v1/order/buyer/buy-now" method:"post" tags:"Order-Buyer" summary:"Create order by buy-now items"`
	pb.CreateOrderBuyNowReq
}

type CreateOrderBuyNowRes = pb.CreateOrderBuyNowRes

type RequestPayReq struct {
	g.Meta `path:"/v1/order/buyer/request-pay" method:"post" tags:"Order-Buyer" summary:"Create pay intent for order"`
	pb.RequestPayReq
}

type RequestPayRes = pb.RequestPayRes

type CancelMyOrderReq struct {
	g.Meta `path:"/v1/order/buyer/cancel" method:"post" tags:"Order-Buyer" summary:"Cancel my order"`
	pb.CancelMyOrderReq
}

type CancelMyOrderRes = pb.CancelMyOrderRes

type GetMyOrderDetailReq struct {
	g.Meta  `path:"/v1/order/buyer/{order_no}" method:"get" tags:"Order-Buyer" summary:"Get my order detail"`
	OrderNo string `json:"order_no" v:"required#order_no is required"`
}

type GetMyOrderDetailRes = pb.GetMyOrderDetailRes

type ListMyOrdersReq struct {
	g.Meta `path:"/v1/order/buyer/list" method:"post" tags:"Order-Buyer" summary:"List my orders with cursor"`
	pb.ListMyOrdersReq
}

type ListMyOrdersRes = pb.ListMyOrdersRes

type ListShopOrdersReq struct {
	g.Meta `path:"/v1/order/seller/list" method:"post" tags:"Order-Seller" summary:"List shop sub orders"`
	pb.ListShopOrdersReq
}

type ListShopOrdersRes = pb.ListShopOrdersRes

type GetShopOrderDetailReq struct {
	g.Meta     `path:"/v1/order/seller/sub/{sub_order_no}" method:"get" tags:"Order-Seller" summary:"Get sub order detail"`
	SubOrderNo string `json:"sub_order_no" v:"required#sub_order_no is required"`
}

type GetShopOrderDetailRes = pb.GetShopOrderDetailRes

type MarkSubOrderShippedReq struct {
	g.Meta `path:"/v1/order/seller/sub/ship" method:"post" tags:"Order-Seller" summary:"Mark sub order shipped"`
	pb.MarkSubOrderShippedReq
}

type MarkSubOrderShippedRes = pb.MarkSubOrderShippedRes

type CloseOrderIfUnpaidReq struct {
	g.Meta `path:"/v1/order/internal/close-unpaid" method:"post" tags:"Order-Internal" summary:"Close unpaid order by timeout worker"`
	pb.CloseOrderIfUnpaidReq
}

type CloseOrderIfUnpaidRes = pb.CloseOrderIfUnpaidRes

type HandlePayCallbackReq struct {
	g.Meta `path:"/v1/order/internal/pay-callback" method:"post" tags:"Order-Internal" summary:"Handle payment callback"`
	pb.HandlePayCallbackReq
}

type HandlePayCallbackRes = pb.HandlePayCallbackRes

type CompleteOrderReq struct {
	g.Meta `path:"/v1/order/internal/complete" method:"post" tags:"Order-Internal" summary:"Complete order and trigger points grant"`
	pb.CompleteOrderReq
}

type CompleteOrderRes = pb.CompleteOrderRes

type GetOrderSnapshotByNoReq struct {
	g.Meta  `path:"/v1/order/internal/{order_no}/snapshot" method:"get" tags:"Order-Internal" summary:"Get order snapshot by order_no"`
	OrderNo string `json:"order_no" v:"required#order_no is required"`
}

type GetOrderSnapshotByNoRes = pb.GetOrderSnapshotByNoRes
