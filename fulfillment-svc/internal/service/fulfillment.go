// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/fulfillment-svc/api/v1"
)

type (
	IFulfillment interface {
		CreateShipment(ctx context.Context, req *v1.CreateShipmentReq) (*v1.CreateShipmentRes, error)
		MarkShipmentShipped(ctx context.Context, req *v1.MarkShipmentShippedReq) (*v1.MarkShipmentShippedRes, error)
		ListShopShipments(ctx context.Context, req *v1.ListShopShipmentsReq) (*v1.ListShopShipmentsRes, error)
		GetShipmentDetail(ctx context.Context, req *v1.GetShipmentDetailReq) (*v1.GetShipmentDetailRes, error)
		GetMyOrderLogistics(ctx context.Context, req *v1.GetMyOrderLogisticsReq) (*v1.GetMyOrderLogisticsRes, error)
		IngestTrackingCallback(ctx context.Context, req *v1.IngestTrackingCallbackReq) (*v1.IngestTrackingCallbackRes, error)
		SyncTrackingByShipment(ctx context.Context, req *v1.SyncTrackingByShipmentReq) (*v1.SyncTrackingByShipmentRes, error)
		GetShipmentSnapshot(ctx context.Context, req *v1.GetShipmentSnapshotReq) (*v1.GetShipmentSnapshotRes, error)
	}
)

var (
	localFulfillment IFulfillment
)

func Fulfillment() IFulfillment {
	if localFulfillment == nil {
		panic("implement not found for interface IFulfillment, forgot register?")
	}
	return localFulfillment
}

func RegisterFulfillment(i IFulfillment) {
	localFulfillment = i
}
