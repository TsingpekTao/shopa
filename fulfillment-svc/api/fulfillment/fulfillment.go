package fulfillment

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/fulfillment-svc/api/fulfillment/v1"
)

type IFulfillmentV1 interface {
	CreateShipment(ctx context.Context, req *v1.CreateShipmentReq) (res *v1.CreateShipmentRes, err error)
	MarkShipmentShipped(ctx context.Context, req *v1.MarkShipmentShippedReq) (res *v1.MarkShipmentShippedRes, err error)
	MarkShipmentShippedAlias(ctx context.Context, req *v1.MarkShipmentShippedAliasReq) (res *v1.MarkShipmentShippedAliasRes, err error)
	ListShopShipments(ctx context.Context, req *v1.ListShopShipmentsReq) (res *v1.ListShopShipmentsRes, err error)
	ListShopShipmentsAlias(ctx context.Context, req *v1.ListShopShipmentsAliasReq) (res *v1.ListShopShipmentsAliasRes, err error)
	GetShipmentDetail(ctx context.Context, req *v1.GetShipmentDetailReq) (res *v1.GetShipmentDetailRes, err error)
	GetMyOrderLogistics(ctx context.Context, req *v1.GetMyOrderLogisticsReq) (res *v1.GetMyOrderLogisticsRes, err error)
	IngestTrackingCallback(ctx context.Context, req *v1.IngestTrackingCallbackReq) (res *v1.IngestTrackingCallbackRes, err error)
	SyncTrackingByShipment(ctx context.Context, req *v1.SyncTrackingByShipmentReq) (res *v1.SyncTrackingByShipmentRes, err error)
	GetShipmentSnapshot(ctx context.Context, req *v1.GetShipmentSnapshotReq) (res *v1.GetShipmentSnapshotRes, err error)
}
