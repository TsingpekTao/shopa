package fulfillment

import (
	"context"

	fulfillmentv1 "github.com/TsingpekTao/shopa/fulfillment-svc/api/fulfillment/v1"
	pb "github.com/TsingpekTao/shopa/fulfillment-svc/api/v1"
)

func (c *ControllerV1) CreateShipment(ctx context.Context, req *fulfillmentv1.CreateShipmentReq) (*fulfillmentv1.CreateShipmentRes, error) {
	return c.fulfillment.CreateShipment(ctx, &req.CreateShipmentReq)
}

func (c *ControllerV1) MarkShipmentShipped(ctx context.Context, req *fulfillmentv1.MarkShipmentShippedReq) (*fulfillmentv1.MarkShipmentShippedRes, error) {
	return c.fulfillment.MarkShipmentShipped(ctx, &req.MarkShipmentShippedReq)
}

func (c *ControllerV1) MarkShipmentShippedAlias(ctx context.Context, req *fulfillmentv1.MarkShipmentShippedAliasReq) (*fulfillmentv1.MarkShipmentShippedAliasRes, error) {
	return c.fulfillment.MarkShipmentShipped(ctx, &req.MarkShipmentShippedReq)
}

func (c *ControllerV1) ListShopShipments(ctx context.Context, req *fulfillmentv1.ListShopShipmentsReq) (*fulfillmentv1.ListShopShipmentsRes, error) {
	return c.fulfillment.ListShopShipments(ctx, &req.ListShopShipmentsReq)
}

func (c *ControllerV1) ListShopShipmentsAlias(ctx context.Context, req *fulfillmentv1.ListShopShipmentsAliasReq) (*fulfillmentv1.ListShopShipmentsAliasRes, error) {
	return c.fulfillment.ListShopShipments(ctx, &pb.ListShopShipmentsReq{
		ShopNo:     req.ShopNo,
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
		Statuses:   req.Statuses,
	})
}

func (c *ControllerV1) GetShipmentDetail(ctx context.Context, req *fulfillmentv1.GetShipmentDetailReq) (*fulfillmentv1.GetShipmentDetailRes, error) {
	return c.fulfillment.GetShipmentDetail(ctx, &pb.GetShipmentDetailReq{ShipmentNo: req.ShipmentNo})
}

func (c *ControllerV1) GetMyOrderLogistics(ctx context.Context, req *fulfillmentv1.GetMyOrderLogisticsReq) (*fulfillmentv1.GetMyOrderLogisticsRes, error) {
	return c.fulfillment.GetMyOrderLogistics(ctx, &pb.GetMyOrderLogisticsReq{OrderNo: req.OrderNo})
}

func (c *ControllerV1) IngestTrackingCallback(ctx context.Context, req *fulfillmentv1.IngestTrackingCallbackReq) (*fulfillmentv1.IngestTrackingCallbackRes, error) {
	return c.fulfillment.IngestTrackingCallback(ctx, &req.IngestTrackingCallbackReq)
}

func (c *ControllerV1) SyncTrackingByShipment(ctx context.Context, req *fulfillmentv1.SyncTrackingByShipmentReq) (*fulfillmentv1.SyncTrackingByShipmentRes, error) {
	return c.fulfillment.SyncTrackingByShipment(ctx, &req.SyncTrackingByShipmentReq)
}

func (c *ControllerV1) GetShipmentSnapshot(ctx context.Context, req *fulfillmentv1.GetShipmentSnapshotReq) (*fulfillmentv1.GetShipmentSnapshotRes, error) {
	return c.fulfillment.GetShipmentSnapshot(ctx, &pb.GetShipmentSnapshotReq{ShipmentNo: req.ShipmentNo})
}
