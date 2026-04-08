package v1

import (
	aftersalev1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type ListRefundBatchesReq struct {
	g.Meta     `path:"/v1/aftersale/seller/shops/{shopNo}/refund-batches" method:"get" tags:"Seller" summary:"List shop refund batches"`
	ShopNo     string                         `json:"shopNo" in:"path" v:"required#shopNo is required"`
	PageSize   int32                          `json:"page_size"`
	NextCursor string                         `json:"next_cursor"`
	Statuses   []aftersalev1.AfterSaleStatus  `json:"statuses"`
}

type ListRefundBatchesRes = aftersalev1.ListShopRefundBatchesRes

type GetRefundBatchDetailReq struct {
	g.Meta        `path:"/v1/aftersale/seller/refund-batches/{refundBatchNo}" method:"get" tags:"Seller" summary:"Get shop refund batch detail"`
	RefundBatchNo string `json:"refundBatchNo" in:"path" v:"required#refundBatchNo is required"`
	ShopNo        string `json:"shopNo" in:"query" v:"required#shopNo is required"`
}

type GetRefundBatchDetailRes = aftersalev1.GetShopRefundBatchDetailRes

type ApproveRefundBatchReq struct {
	g.Meta `path:"/v1/aftersale/seller/refund-batches:approve" method:"post" tags:"Seller" summary:"Approve refund batch"`
	aftersalev1.ApproveRefundBatchReq
}

type ApproveRefundBatchRes = aftersalev1.ApproveRefundBatchRes

type RejectRefundBatchReq struct {
	g.Meta `path:"/v1/aftersale/seller/refund-batches:reject" method:"post" tags:"Seller" summary:"Reject refund batch"`
	aftersalev1.RejectRefundBatchReq
}

type RejectRefundBatchRes = aftersalev1.RejectRefundBatchRes
