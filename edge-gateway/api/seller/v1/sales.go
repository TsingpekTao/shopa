package v1

import "github.com/gogf/gf/v2/frame/g"

type GetSalesAnalyticsReq struct {
	g.Meta `path:"/v1/seller/sales/analytics" method:"get" tags:"Seller" summary:"Get seller sales analytics aggregated by gateway"`
	ShopNo string `json:"shopNo" in:"query"`
	Range  string `json:"range" in:"query" v:"required|in:7D,14D,30D,8W#range is required|range must be one of 7D,14D,30D,8W"`
}

type SellerSalesSummary struct {
	Gmv            uint64  `json:"gmv"`
	PaidOrderCount uint64  `json:"paidOrderCount"`
	PaidBuyerCount uint64  `json:"paidBuyerCount"`
	RefundAmount   uint64  `json:"refundAmount"`
	RefundRate     float64 `json:"refundRate"`
	AvgOrderValue  uint64  `json:"avgOrderValue"`
}

type SellerSalesBucket struct {
	BucketKey      string  `json:"bucketKey"`
	BucketLabel    string  `json:"bucketLabel"`
	Gmv            uint64  `json:"gmv"`
	PaidOrderCount uint64  `json:"paidOrderCount"`
	PaidBuyerCount uint64  `json:"paidBuyerCount"`
	RefundAmount   uint64  `json:"refundAmount"`
	RefundRate     float64 `json:"refundRate"`
}

type GetSalesAnalyticsRes struct {
	Summary        SellerSalesSummary  `json:"summary"`
	Series         []SellerSalesBucket `json:"series"`
	Partial        bool                `json:"partial"`
	DegradedFields []string            `json:"degradedFields"`
}
