package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type GetWorkbenchReq struct {
	g.Meta `path:"/v1/seller/workbench" method:"get" tags:"Seller" summary:"Get seller workbench overview aggregated by gateway"`
}

type SellerShopSummary struct {
	ShopNo          string `json:"shopNo"`
	ShopName        string `json:"shopName"`
	ShopDisplayName string `json:"shopDisplayName"`
	ShopStatusCode  string `json:"shopStatusCode"`
	BuyerVisible    bool   `json:"buyerVisible"`
	UpdatedAt       string `json:"updatedAt"`
}

type SellerApplicationSummary struct {
	ApplicationNo         string `json:"applicationNo"`
	ApplicationStatusCode string `json:"applicationStatusCode"`
	Version               int32  `json:"version"`
	SubmittedAt           string `json:"submittedAt"`
	UpdatedAt             string `json:"updatedAt"`
}

type SellerProductSummary struct {
	Total     uint64 `json:"total"`
	OnShelf   uint64 `json:"onShelf"`
	OffShelf  uint64 `json:"offShelf"`
	Reviewing uint64 `json:"reviewing"`
	Draft     uint64 `json:"draft"`
	Rejected  uint64 `json:"rejected"`
}

type GetWorkbenchRes struct {
	UserID uint64 `json:"userId"`

	// Workbench only returns the first 10 shops.
	Shops          []SellerShopSummary `json:"shops"`
	ShopsTotal     uint64              `json:"shopsTotal"`
	ShopsTruncated bool                `json:"shopsTruncated"`

	// Workbench only returns the latest 5 applications.
	LatestApplications    []SellerApplicationSummary `json:"latestApplications"`
	ApplicationsTotal     uint64                     `json:"applicationsTotal"`
	ApplicationsTruncated bool                       `json:"applicationsTruncated"`

	ProductSummary SellerProductSummary `json:"productSummary"`
	Partial        bool                 `json:"partial"`
	DegradedFields []string             `json:"degradedFields"`
}

type GetShopDashboardReq struct {
	g.Meta `path:"/v1/seller/shops/{shopNo}/dashboard" method:"get" tags:"Seller" summary:"Get seller shop dashboard aggregated by gateway"`
	ShopNo string `json:"shopNo" in:"path" v:"required"`
}

type InventoryRiskSummary struct {
	LowStockSkuCount   uint64 `json:"lowStockSkuCount"`
	OutOfStockSkuCount uint64 `json:"outOfStockSkuCount"`
}

type GetShopDashboardRes struct {
	ShopNo         string               `json:"shopNo"`
	ShopName       string               `json:"shopName"`
	ShopStatusCode string               `json:"shopStatusCode"`
	ProductSummary SellerProductSummary `json:"productSummary"`
	InventoryRisk  InventoryRiskSummary `json:"inventoryRisk"`
	Partial        bool                 `json:"partial"`
	DegradedFields []string             `json:"degradedFields"`
}
