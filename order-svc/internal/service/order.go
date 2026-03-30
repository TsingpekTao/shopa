// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
)

type (
	IOrder interface {
		CreateOrderFromCart(ctx context.Context, req *v1.CreateOrderFromCartReq) (*v1.CreateOrderFromCartRes, error)
		CreateOrderBuyNow(ctx context.Context, req *v1.CreateOrderBuyNowReq) (*v1.CreateOrderBuyNowRes, error)
		RequestPay(ctx context.Context, req *v1.RequestPayReq) (*v1.RequestPayRes, error)
		CancelMyOrder(ctx context.Context, req *v1.CancelMyOrderReq) (*v1.CancelMyOrderRes, error)
		GetMyOrderDetail(ctx context.Context, req *v1.GetMyOrderDetailReq) (*v1.GetMyOrderDetailRes, error)
		ListMyOrders(ctx context.Context, req *v1.ListMyOrdersReq) (*v1.ListMyOrdersRes, error)
		ListShopOrders(ctx context.Context, req *v1.ListShopOrdersReq) (*v1.ListShopOrdersRes, error)
		GetShopOrderDetail(ctx context.Context, req *v1.GetShopOrderDetailReq) (*v1.GetShopOrderDetailRes, error)
		MarkSubOrderShipped(ctx context.Context, req *v1.MarkSubOrderShippedReq) (*v1.MarkSubOrderShippedRes, error)
		CloseOrderIfUnpaid(ctx context.Context, req *v1.CloseOrderIfUnpaidReq) (*v1.CloseOrderIfUnpaidRes, error)
		HandlePayCallback(ctx context.Context, req *v1.HandlePayCallbackReq) (*v1.HandlePayCallbackRes, error)
		GetOrderSnapshotByNo(ctx context.Context, req *v1.GetOrderSnapshotByNoReq) (*v1.GetOrderSnapshotByNoRes, error)
	}
)

var (
	localOrder IOrder
)

func Order() IOrder {
	if localOrder == nil {
		panic("implement not found for interface IOrder, forgot register?")
	}
	return localOrder
}

func RegisterOrder(i IOrder) {
	localOrder = i
}
