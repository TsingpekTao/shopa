package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/fulfillment-svc/api/v1"
	"github.com/TsingpekTao/shopa/fulfillment-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedSellerFulfillmentServiceServer
	v1.UnimplementedBuyerFulfillmentServiceServer
	v1.UnimplementedInternalFulfillmentServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterSellerFulfillmentServiceServer(s.Server, ctrl)
	v1.RegisterBuyerFulfillmentServiceServer(s.Server, ctrl)
	v1.RegisterInternalFulfillmentServiceServer(s.Server, ctrl)
}

func (*Controller) CreateShipment(ctx context.Context, req *v1.CreateShipmentReq) (res *v1.CreateShipmentRes, err error) {
	return service.Fulfillment().CreateShipment(ctx, req)
}

func (*Controller) MarkShipmentShipped(ctx context.Context, req *v1.MarkShipmentShippedReq) (res *v1.MarkShipmentShippedRes, err error) {
	return service.Fulfillment().MarkShipmentShipped(ctx, req)
}

func (*Controller) ListShopShipments(ctx context.Context, req *v1.ListShopShipmentsReq) (res *v1.ListShopShipmentsRes, err error) {
	return service.Fulfillment().ListShopShipments(ctx, req)
}

func (*Controller) GetShipmentDetail(ctx context.Context, req *v1.GetShipmentDetailReq) (res *v1.GetShipmentDetailRes, err error) {
	return service.Fulfillment().GetShipmentDetail(ctx, req)
}

func (*Controller) GetMyOrderLogistics(ctx context.Context, req *v1.GetMyOrderLogisticsReq) (res *v1.GetMyOrderLogisticsRes, err error) {
	return service.Fulfillment().GetMyOrderLogistics(ctx, req)
}

func (*Controller) IngestTrackingCallback(ctx context.Context, req *v1.IngestTrackingCallbackReq) (res *v1.IngestTrackingCallbackRes, err error) {
	return service.Fulfillment().IngestTrackingCallback(ctx, req)
}

func (*Controller) SyncTrackingByShipment(ctx context.Context, req *v1.SyncTrackingByShipmentReq) (res *v1.SyncTrackingByShipmentRes, err error) {
	return service.Fulfillment().SyncTrackingByShipment(ctx, req)
}

func (*Controller) GetShipmentSnapshot(ctx context.Context, req *v1.GetShipmentSnapshotReq) (res *v1.GetShipmentSnapshotRes, err error) {
	return service.Fulfillment().GetShipmentSnapshot(ctx, req)
}
