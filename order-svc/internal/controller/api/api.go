package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/TsingpekTao/shopa/order-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedBuyerOrderServiceServer
	v1.UnimplementedSellerOrderServiceServer
	v1.UnimplementedInternalOrderServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterBuyerOrderServiceServer(s.Server, ctrl)
	v1.RegisterSellerOrderServiceServer(s.Server, ctrl)
	v1.RegisterInternalOrderServiceServer(s.Server, ctrl)
}

func (*Controller) CreateOrderFromCart(ctx context.Context, req *v1.CreateOrderFromCartReq) (*v1.CreateOrderFromCartRes, error) {
	return service.Order().CreateOrderFromCart(ctx, req)
}

func (*Controller) CreateOrderBuyNow(ctx context.Context, req *v1.CreateOrderBuyNowReq) (*v1.CreateOrderBuyNowRes, error) {
	return service.Order().CreateOrderBuyNow(ctx, req)
}

func (*Controller) RequestPay(ctx context.Context, req *v1.RequestPayReq) (*v1.RequestPayRes, error) {
	return service.Order().RequestPay(ctx, req)
}

func (*Controller) CancelMyOrder(ctx context.Context, req *v1.CancelMyOrderReq) (*v1.CancelMyOrderRes, error) {
	return service.Order().CancelMyOrder(ctx, req)
}

func (*Controller) GetMyOrderDetail(ctx context.Context, req *v1.GetMyOrderDetailReq) (*v1.GetMyOrderDetailRes, error) {
	return service.Order().GetMyOrderDetail(ctx, req)
}

func (*Controller) ListMyOrders(ctx context.Context, req *v1.ListMyOrdersReq) (*v1.ListMyOrdersRes, error) {
	return service.Order().ListMyOrders(ctx, req)
}

func (*Controller) ListShopOrders(ctx context.Context, req *v1.ListShopOrdersReq) (*v1.ListShopOrdersRes, error) {
	return service.Order().ListShopOrders(ctx, req)
}

func (*Controller) GetShopOrderDetail(ctx context.Context, req *v1.GetShopOrderDetailReq) (*v1.GetShopOrderDetailRes, error) {
	return service.Order().GetShopOrderDetail(ctx, req)
}

func (*Controller) MarkSubOrderShipped(ctx context.Context, req *v1.MarkSubOrderShippedReq) (*v1.MarkSubOrderShippedRes, error) {
	return service.Order().MarkSubOrderShipped(ctx, req)
}

func (*Controller) CloseOrderIfUnpaid(ctx context.Context, req *v1.CloseOrderIfUnpaidReq) (*v1.CloseOrderIfUnpaidRes, error) {
	return service.Order().CloseOrderIfUnpaid(ctx, req)
}

func (*Controller) HandlePayCallback(ctx context.Context, req *v1.HandlePayCallbackReq) (*v1.HandlePayCallbackRes, error) {
	return service.Order().HandlePayCallback(ctx, req)
}

func (*Controller) CompleteOrder(ctx context.Context, req *v1.CompleteOrderReq) (*v1.CompleteOrderRes, error) {
	return service.Order().CompleteOrder(ctx, req)
}

func (*Controller) GetOrderSnapshotByNo(ctx context.Context, req *v1.GetOrderSnapshotByNoReq) (*v1.GetOrderSnapshotByNoRes, error) {
	return service.Order().GetOrderSnapshotByNo(ctx, req)
}
