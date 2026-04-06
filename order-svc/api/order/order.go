package order

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/order-svc/api/order/v1"
)

type IOrderV1 interface {
	CreateOrderFromCart(ctx context.Context, req *v1.CreateOrderFromCartReq) (res *v1.CreateOrderFromCartRes, err error)
	CreateOrderBuyNow(ctx context.Context, req *v1.CreateOrderBuyNowReq) (res *v1.CreateOrderBuyNowRes, err error)
	RequestPay(ctx context.Context, req *v1.RequestPayReq) (res *v1.RequestPayRes, err error)
	CancelMyOrder(ctx context.Context, req *v1.CancelMyOrderReq) (res *v1.CancelMyOrderRes, err error)
	GetMyOrderDetail(ctx context.Context, req *v1.GetMyOrderDetailReq) (res *v1.GetMyOrderDetailRes, err error)
	ListMyOrders(ctx context.Context, req *v1.ListMyOrdersReq) (res *v1.ListMyOrdersRes, err error)
	ListShopOrders(ctx context.Context, req *v1.ListShopOrdersReq) (res *v1.ListShopOrdersRes, err error)
	GetShopOrderDetail(ctx context.Context, req *v1.GetShopOrderDetailReq) (res *v1.GetShopOrderDetailRes, err error)
	MarkSubOrderShipped(ctx context.Context, req *v1.MarkSubOrderShippedReq) (res *v1.MarkSubOrderShippedRes, err error)
	CloseOrderIfUnpaid(ctx context.Context, req *v1.CloseOrderIfUnpaidReq) (res *v1.CloseOrderIfUnpaidRes, err error)
	HandlePayCallback(ctx context.Context, req *v1.HandlePayCallbackReq) (res *v1.HandlePayCallbackRes, err error)
	CompleteOrder(ctx context.Context, req *v1.CompleteOrderReq) (res *v1.CompleteOrderRes, err error)
	GetOrderSnapshotByNo(ctx context.Context, req *v1.GetOrderSnapshotByNoReq) (res *v1.GetOrderSnapshotByNoRes, err error)
}
