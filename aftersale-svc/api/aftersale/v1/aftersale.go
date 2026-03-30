package v1

import (
	pb "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreateAfterSaleReq struct {
	g.Meta `path:"/v1/aftersale/buyer/create" method:"post" tags:"AfterSale-Buyer" summary:"Create after-sale case"`
	pb.CreateAfterSaleReq
}

type CreateAfterSaleRes = pb.CreateAfterSaleRes

type CreateAfterSaleAliasReq struct {
	g.Meta `path:"/v1/aftersale/buyer/cases:create" method:"post" tags:"AfterSale-Buyer" summary:"Create after-sale case (alias)"`
	pb.CreateAfterSaleReq
}

type CreateAfterSaleAliasRes = pb.CreateAfterSaleRes

type CancelAfterSaleReq struct {
	g.Meta `path:"/v1/aftersale/buyer/cancel" method:"post" tags:"AfterSale-Buyer" summary:"Cancel after-sale case"`
	pb.CancelAfterSaleReq
}

type CancelAfterSaleRes = pb.CancelAfterSaleRes

type CancelAfterSaleAliasReq struct {
	g.Meta `path:"/v1/aftersale/buyer/cases:cancel" method:"post" tags:"AfterSale-Buyer" summary:"Cancel after-sale case (alias)"`
	pb.CancelAfterSaleReq
}

type CancelAfterSaleAliasRes = pb.CancelAfterSaleRes

type GetMyAfterSaleDetailReq struct {
	g.Meta      `path:"/v1/aftersale/buyer/{after_sale_no}" method:"get" tags:"AfterSale-Buyer" summary:"Get my after-sale detail"`
	AfterSaleNo string `json:"after_sale_no" v:"required#after_sale_no is required"`
}

type GetMyAfterSaleDetailRes = pb.GetMyAfterSaleDetailRes

type GetMyAfterSaleDetailAliasReq struct {
	g.Meta      `path:"/v1/aftersale/buyer/cases/{after_sale_no}" method:"get" tags:"AfterSale-Buyer" summary:"Get my after-sale detail (alias)"`
	AfterSaleNo string `json:"after_sale_no" v:"required#after_sale_no is required"`
}

type GetMyAfterSaleDetailAliasRes = pb.GetMyAfterSaleDetailRes

type ListMyAfterSalesReq struct {
	g.Meta `path:"/v1/aftersale/buyer/list" method:"post" tags:"AfterSale-Buyer" summary:"List my after-sales"`
	pb.ListMyAfterSalesReq
}

type ListMyAfterSalesRes = pb.ListMyAfterSalesRes

type ListMyAfterSalesAliasReq struct {
	g.Meta     `path:"/v1/aftersale/buyer/cases" method:"get" tags:"AfterSale-Buyer" summary:"List my after-sales (alias)"`
	PageSize   int32                `json:"page_size"`
	NextCursor string               `json:"next_cursor"`
	Statuses   []pb.AfterSaleStatus `json:"statuses"`
}

type ListMyAfterSalesAliasRes = pb.ListMyAfterSalesRes

type ListShopAfterSalesReq struct {
	g.Meta `path:"/v1/aftersale/seller/list" method:"post" tags:"AfterSale-Seller" summary:"List shop after-sales"`
	pb.ListShopAfterSalesReq
}

type ListShopAfterSalesRes = pb.ListShopAfterSalesRes

type ListShopAfterSalesAliasReq struct {
	g.Meta     `path:"/v1/aftersale/seller/shops/{shop_no}/cases" method:"get" tags:"AfterSale-Seller" summary:"List shop after-sales (alias)"`
	ShopNo     string               `json:"shop_no" v:"required#shop_no is required"`
	PageSize   int32                `json:"page_size"`
	NextCursor string               `json:"next_cursor"`
	Statuses   []pb.AfterSaleStatus `json:"statuses"`
}

type ListShopAfterSalesAliasRes = pb.ListShopAfterSalesRes

type GetShopAfterSaleDetailReq struct {
	g.Meta      `path:"/v1/aftersale/seller/{after_sale_no}" method:"get" tags:"AfterSale-Seller" summary:"Get shop after-sale detail"`
	AfterSaleNo string `json:"after_sale_no" v:"required#after_sale_no is required"`
}

type GetShopAfterSaleDetailRes = pb.GetShopAfterSaleDetailRes

type GetShopAfterSaleDetailAliasReq struct {
	g.Meta      `path:"/v1/aftersale/seller/cases/{after_sale_no}" method:"get" tags:"AfterSale-Seller" summary:"Get shop after-sale detail (alias)"`
	AfterSaleNo string `json:"after_sale_no" v:"required#after_sale_no is required"`
}

type GetShopAfterSaleDetailAliasRes = pb.GetShopAfterSaleDetailRes

type ApproveAfterSaleReq struct {
	g.Meta `path:"/v1/aftersale/seller/approve" method:"post" tags:"AfterSale-Seller" summary:"Approve after-sale case"`
	pb.ApproveAfterSaleReq
}

type ApproveAfterSaleRes = pb.ApproveAfterSaleRes

type ApproveAfterSaleAliasReq struct {
	g.Meta `path:"/v1/aftersale/seller/cases:approve" method:"post" tags:"AfterSale-Seller" summary:"Approve after-sale case (alias)"`
	pb.ApproveAfterSaleReq
}

type ApproveAfterSaleAliasRes = pb.ApproveAfterSaleRes

type RejectAfterSaleReq struct {
	g.Meta `path:"/v1/aftersale/seller/reject" method:"post" tags:"AfterSale-Seller" summary:"Reject after-sale case"`
	pb.RejectAfterSaleReq
}

type RejectAfterSaleRes = pb.RejectAfterSaleRes

type RejectAfterSaleAliasReq struct {
	g.Meta `path:"/v1/aftersale/seller/cases:reject" method:"post" tags:"AfterSale-Seller" summary:"Reject after-sale case (alias)"`
	pb.RejectAfterSaleReq
}

type RejectAfterSaleAliasRes = pb.RejectAfterSaleRes

type ExecuteRefundTaskReq struct {
	g.Meta `path:"/v1/aftersale/internal/refund/execute" method:"post" tags:"AfterSale-Internal" summary:"Execute refund task"`
	pb.ExecuteRefundTaskReq
}

type ExecuteRefundTaskRes = pb.ExecuteRefundTaskRes

type RetryRefundTaskReq struct {
	g.Meta `path:"/v1/aftersale/internal/refund/retry" method:"post" tags:"AfterSale-Internal" summary:"Retry refund task"`
	pb.RetryRefundTaskReq
}

type RetryRefundTaskRes = pb.RetryRefundTaskRes

type GetAfterSaleSnapshotByNoReq struct {
	g.Meta      `path:"/v1/aftersale/internal/{after_sale_no}/snapshot" method:"get" tags:"AfterSale-Internal" summary:"Get after-sale snapshot"`
	AfterSaleNo string `json:"after_sale_no" v:"required#after_sale_no is required"`
}

type GetAfterSaleSnapshotByNoRes = pb.GetAfterSaleSnapshotByNoRes
