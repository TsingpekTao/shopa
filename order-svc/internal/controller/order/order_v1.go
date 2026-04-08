package order

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/order-svc/api/order/v1"
	pb "github.com/TsingpekTao/shopa/order-svc/api/v1"
)

func (c *ControllerV1) CreateOrderFromCart(ctx context.Context, req *httpv1.CreateOrderFromCartReq) (*httpv1.CreateOrderFromCartRes, error) {
	return c.order.CreateOrderFromCart(ctx, &req.CreateOrderFromCartReq)
}

func (c *ControllerV1) CreateOrderBuyNow(ctx context.Context, req *httpv1.CreateOrderBuyNowReq) (*httpv1.CreateOrderBuyNowRes, error) {
	return c.order.CreateOrderBuyNow(ctx, &req.CreateOrderBuyNowReq)
}

func (c *ControllerV1) UpdateMyOrderAddress(ctx context.Context, req *httpv1.UpdateMyOrderAddressReq) (*httpv1.UpdateMyOrderAddressRes, error) {
	return c.order.UpdateMyOrderAddress(ctx, &pb.UpdateMyOrderAddressReq{
		OrderNo:   req.OrderNo,
		AddressId: req.AddressId,
	})
}

func (c *ControllerV1) RequestPay(ctx context.Context, req *httpv1.RequestPayReq) (*httpv1.RequestPayRes, error) {
	return c.order.RequestPay(ctx, &req.RequestPayReq)
}

func (c *ControllerV1) CancelMyOrder(ctx context.Context, req *httpv1.CancelMyOrderReq) (*httpv1.CancelMyOrderRes, error) {
	return c.order.CancelMyOrder(ctx, &req.CancelMyOrderReq)
}

func (c *ControllerV1) ConfirmMyOrderReceived(ctx context.Context, req *httpv1.ConfirmMyOrderReceivedReq) (*httpv1.ConfirmMyOrderReceivedRes, error) {
	return c.order.ConfirmMyOrderReceived(ctx, &req.CompleteOrderReq)
}

func (c *ControllerV1) GetMyOrderDetail(ctx context.Context, req *httpv1.GetMyOrderDetailReq) (*httpv1.GetMyOrderDetailRes, error) {
	return c.order.GetMyOrderDetail(ctx, &pb.GetMyOrderDetailReq{OrderNo: req.OrderNo})
}

func (c *ControllerV1) ListMyOrders(ctx context.Context, req *httpv1.ListMyOrdersReq) (*httpv1.ListMyOrdersRes, error) {
	return c.order.ListMyOrders(ctx, &req.ListMyOrdersReq)
}

func (c *ControllerV1) ListShopOrders(ctx context.Context, req *httpv1.ListShopOrdersReq) (*httpv1.ListShopOrdersRes, error) {
	return c.order.ListShopOrders(ctx, &req.ListShopOrdersReq)
}

func (c *ControllerV1) GetShopOrderDetail(ctx context.Context, req *httpv1.GetShopOrderDetailReq) (*httpv1.GetShopOrderDetailRes, error) {
	return c.order.GetShopOrderDetail(ctx, &pb.GetShopOrderDetailReq{SubOrderNo: req.SubOrderNo})
}

func (c *ControllerV1) MarkSubOrderShipped(ctx context.Context, req *httpv1.MarkSubOrderShippedReq) (*httpv1.MarkSubOrderShippedRes, error) {
	return c.order.MarkSubOrderShipped(ctx, &req.MarkSubOrderShippedReq)
}

func (c *ControllerV1) CloseOrderIfUnpaid(ctx context.Context, req *httpv1.CloseOrderIfUnpaidReq) (*httpv1.CloseOrderIfUnpaidRes, error) {
	return c.order.CloseOrderIfUnpaid(ctx, &req.CloseOrderIfUnpaidReq)
}

func (c *ControllerV1) HandlePayCallback(ctx context.Context, req *httpv1.HandlePayCallbackReq) (*httpv1.HandlePayCallbackRes, error) {
	return c.order.HandlePayCallback(ctx, &req.HandlePayCallbackReq)
}

func (c *ControllerV1) CompleteOrder(ctx context.Context, req *httpv1.CompleteOrderReq) (*httpv1.CompleteOrderRes, error) {
	return c.order.CompleteOrder(ctx, &req.CompleteOrderReq)
}

func (c *ControllerV1) GetOrderSnapshotByNo(ctx context.Context, req *httpv1.GetOrderSnapshotByNoReq) (*httpv1.GetOrderSnapshotByNoRes, error) {
	return c.order.GetOrderSnapshotByNo(ctx, &pb.GetOrderSnapshotByNoReq{OrderNo: req.OrderNo})
}
