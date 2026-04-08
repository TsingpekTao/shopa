package v1

import (
	aftersalev1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	orderv1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type PreviewRefundReq struct {
	g.Meta          `path:"/v1/aftersale/buyer/refunds:preview" method:"post" tags:"Me" summary:"Preview pre-shipment refund"`
	OrderNo         string   `json:"orderNo" v:"required#orderNo is required"`
	SubOrderNo      string   `json:"subOrderNo" v:"required#subOrderNo is required"`
	ItemNo          string   `json:"itemNo"`
	SelectedItemNos []string `json:"selectedItemNos"`
}

type PreviewRefundRes struct {
	Snapshot   *orderv1.RefundSubOrderSnapshot `json:"snapshot"`
	ScopeCode  string                          `json:"scopeCode"`
	ScopeLabel string                          `json:"scopeLabel"`
}

type ApplyRefundBatchReq struct {
	g.Meta `path:"/v1/aftersale/buyer/refunds:apply" method:"post" tags:"Me" summary:"Apply refund batch"`
	aftersalev1.ApplyRefundBatchReq
}

type ApplyRefundBatchRes = aftersalev1.ApplyRefundBatchRes

type ListRefundBatchesReq struct {
	g.Meta     `path:"/v1/aftersale/buyer/refund-batches" method:"get" tags:"Me" summary:"List my refund batches"`
	PageSize   int32                           `json:"page_size"`
	NextCursor string                          `json:"next_cursor"`
	Statuses   []aftersalev1.AfterSaleStatus   `json:"statuses"`
}

type ListRefundBatchesRes = aftersalev1.ListMyRefundBatchesRes

type GetRefundBatchDetailReq struct {
	g.Meta        `path:"/v1/aftersale/buyer/refund-batches/{refundBatchNo}" method:"get" tags:"Me" summary:"Get my refund batch detail"`
	RefundBatchNo string `json:"refundBatchNo" in:"path" v:"required#refundBatchNo is required"`
}

type GetRefundBatchDetailRes = aftersalev1.GetMyRefundBatchDetailRes

type CancelRefundBatchReq struct {
	g.Meta `path:"/v1/aftersale/buyer/refund-batches:cancel" method:"post" tags:"Me" summary:"Cancel refund batch"`
	aftersalev1.CancelRefundBatchReq
}

type CancelRefundBatchRes = aftersalev1.CancelRefundBatchRes
