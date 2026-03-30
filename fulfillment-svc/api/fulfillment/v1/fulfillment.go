package v1

import (
	pb "github.com/TsingpekTao/shopa/fulfillment-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreateShipmentReq struct {
	g.Meta `path:"/v1/fulfillment/seller/shipments:create" method:"post" tags:"Fulfillment-Seller" summary:"Create shipment"`
	pb.CreateShipmentReq
}

type CreateShipmentRes = pb.CreateShipmentRes

type MarkShipmentShippedReq struct {
	g.Meta `path:"/v1/fulfillment/seller/shipments:mark-shipped" method:"post" tags:"Fulfillment-Seller" summary:"Mark shipment shipped"`
	pb.MarkShipmentShippedReq
}

type MarkShipmentShippedRes = pb.MarkShipmentShippedRes

type MarkShipmentShippedAliasReq struct {
	g.Meta `path:"/v1/fulfillment/seller/shipments:ship" method:"post" tags:"Fulfillment-Seller" summary:"Mark shipment shipped (alias)"`
	pb.MarkShipmentShippedReq
}

type MarkShipmentShippedAliasRes = pb.MarkShipmentShippedRes

type ListShopShipmentsReq struct {
	g.Meta `path:"/v1/fulfillment/seller/shipments:list" method:"post" tags:"Fulfillment-Seller" summary:"List shop shipments"`
	pb.ListShopShipmentsReq
}

type ListShopShipmentsRes = pb.ListShopShipmentsRes

type ListShopShipmentsAliasReq struct {
	g.Meta     `path:"/v1/fulfillment/seller/shops/{shop_no}/shipments" method:"get" tags:"Fulfillment-Seller" summary:"List shop shipments (alias)"`
	ShopNo     string              `json:"shop_no" v:"required#shop_no is required"`
	PageSize   int32               `json:"page_size"`
	NextCursor string              `json:"next_cursor"`
	Statuses   []pb.ShipmentStatus `json:"statuses"`
}

type ListShopShipmentsAliasRes = pb.ListShopShipmentsRes

type GetShipmentDetailReq struct {
	g.Meta     `path:"/v1/fulfillment/seller/shipments/{shipment_no}" method:"get" tags:"Fulfillment-Seller" summary:"Get shipment detail"`
	ShipmentNo string `json:"shipment_no" v:"required#shipment_no is required"`
}

type GetShipmentDetailRes = pb.GetShipmentDetailRes

type GetMyOrderLogisticsReq struct {
	g.Meta  `path:"/v1/fulfillment/buyer/orders/{order_no}/logistics" method:"get" tags:"Fulfillment-Buyer" summary:"Get my order logistics"`
	OrderNo string `json:"order_no" v:"required#order_no is required"`
}

type GetMyOrderLogisticsRes = pb.GetMyOrderLogisticsRes

type IngestTrackingCallbackReq struct {
	g.Meta `path:"/v1/fulfillment/internal/tracking:ingest" method:"post" tags:"Fulfillment-Internal" summary:"Ingest tracking callback"`
	pb.IngestTrackingCallbackReq
}

type IngestTrackingCallbackRes = pb.IngestTrackingCallbackRes

type SyncTrackingByShipmentReq struct {
	g.Meta `path:"/v1/fulfillment/internal/tracking:sync" method:"post" tags:"Fulfillment-Internal" summary:"Sync tracking by shipment"`
	pb.SyncTrackingByShipmentReq
}

type SyncTrackingByShipmentRes = pb.SyncTrackingByShipmentRes

type GetShipmentSnapshotReq struct {
	g.Meta     `path:"/v1/fulfillment/internal/shipments/{shipment_no}/snapshot" method:"get" tags:"Fulfillment-Internal" summary:"Get shipment snapshot"`
	ShipmentNo string `json:"shipment_no" v:"required#shipment_no is required"`
}

type GetShipmentSnapshotRes = pb.GetShipmentSnapshotRes
