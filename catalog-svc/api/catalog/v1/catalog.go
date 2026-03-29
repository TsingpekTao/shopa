package v1

import (
	pb "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreateProductDraftReq struct {
	g.Meta `path:"/v1/catalog/seller/products/draft" method:"post" tags:"Catalog-Seller" summary:"Create product draft"`
	pb.CreateProductDraftReq
}
type CreateProductDraftRes = pb.CreateProductDraftRes

type UpdateProductDraftReq struct {
	g.Meta `path:"/v1/catalog/seller/products/draft" method:"put" tags:"Catalog-Seller" summary:"Update product draft"`
	pb.UpdateProductDraftReq
}
type UpdateProductDraftRes = pb.UpdateProductDraftRes

type UpsertSkuDraftsReq struct {
	g.Meta `path:"/v1/catalog/seller/products/skus:upsert" method:"post" tags:"Catalog-Seller" summary:"Upsert sku drafts"`
	pb.UpsertSkuDraftsReq
}
type UpsertSkuDraftsRes = pb.UpsertSkuDraftsRes

type SubmitProductReviewReq struct {
	g.Meta `path:"/v1/catalog/seller/products/review:submit" method:"post" tags:"Catalog-Seller" summary:"Submit product review"`
	pb.SubmitProductReviewReq
}
type SubmitProductReviewRes = pb.SubmitProductReviewRes

type ResubmitProductReviewReq struct {
	g.Meta `path:"/v1/catalog/seller/products/review:resubmit" method:"post" tags:"Catalog-Seller" summary:"Resubmit product review"`
	pb.ResubmitProductReviewReq
}
type ResubmitProductReviewRes = pb.ResubmitProductReviewRes

type SetProductOnShelfReq struct {
	g.Meta `path:"/v1/catalog/seller/products/on-shelf" method:"post" tags:"Catalog-Seller" summary:"Set product on shelf"`
	pb.SetProductOnShelfReq
}
type SetProductOnShelfRes = pb.SetProductOnShelfRes

type SetProductOffShelfReq struct {
	g.Meta `path:"/v1/catalog/seller/products/off-shelf" method:"post" tags:"Catalog-Seller" summary:"Set product off shelf"`
	pb.SetProductOffShelfReq
}
type SetProductOffShelfRes = pb.SetProductOffShelfRes

type DeleteProductDraftReq struct {
	g.Meta `path:"/v1/catalog/seller/products/draft:delete" method:"post" tags:"Catalog-Seller" summary:"Delete product draft"`
	pb.DeleteProductDraftReq
}
type DeleteProductDraftRes struct {
	Success bool `json:"success"`
}

type GetMyProductReq struct {
	g.Meta `path:"/v1/catalog/seller/products/{spu_no}" method:"get" tags:"Catalog-Seller" summary:"Get my product detail"`
	SpuNo  string `json:"spu_no" dc:"spu no"`
}
type GetMyProductRes = pb.GetMyProductRes

type ListMyProductsReq struct {
	g.Meta   `path:"/v1/catalog/seller/products" method:"get" tags:"Catalog-Seller" summary:"List my products"`
	Page     int32          `json:"page"`
	PageSize int32          `json:"page_size"`
	Statuses []pb.SpuStatus `json:"statuses"`
	Keyword  string         `json:"keyword"`
}
type ListMyProductsRes = pb.ListMyProductsRes

type ListReviewTasksReq struct {
	g.Meta   `path:"/v1/catalog/admin/review/tasks" method:"get" tags:"Catalog-Admin" summary:"List review tasks"`
	Page     int32          `json:"page"`
	PageSize int32          `json:"page_size"`
	Statuses []pb.SpuStatus `json:"statuses"`
	Keyword  string         `json:"keyword"`
}
type ListReviewTasksRes = pb.ListReviewTasksRes

type GetReviewDetailReq struct {
	g.Meta `path:"/v1/catalog/admin/review/{spu_no}" method:"get" tags:"Catalog-Admin" summary:"Get review detail"`
	SpuNo  string `json:"spu_no"`
}
type GetReviewDetailRes = pb.GetReviewDetailRes

type ApproveProductReq struct {
	g.Meta `path:"/v1/catalog/admin/review/approve" method:"post" tags:"Catalog-Admin" summary:"Approve product"`
	pb.ApproveProductReq
}
type ApproveProductRes = pb.ApproveProductRes

type RejectProductReq struct {
	g.Meta `path:"/v1/catalog/admin/review/reject" method:"post" tags:"Catalog-Admin" summary:"Reject product"`
	pb.RejectProductReq
}
type RejectProductRes = pb.RejectProductRes

type FreezeProductReq struct {
	g.Meta `path:"/v1/catalog/admin/review/freeze" method:"post" tags:"Catalog-Admin" summary:"Freeze product"`
	pb.FreezeProductReq
}
type FreezeProductRes = pb.FreezeProductRes

type UnfreezeProductReq struct {
	g.Meta `path:"/v1/catalog/admin/review/unfreeze" method:"post" tags:"Catalog-Admin" summary:"Unfreeze product"`
	pb.UnfreezeProductReq
}
type UnfreezeProductRes = pb.UnfreezeProductRes

type ForceOffShelfReq struct {
	g.Meta `path:"/v1/catalog/admin/review/force-off-shelf" method:"post" tags:"Catalog-Admin" summary:"Force off shelf"`
	pb.ForceOffShelfReq
}
type ForceOffShelfRes = pb.ForceOffShelfRes

type GetProductDetailReq struct {
	g.Meta `path:"/v1/catalog/buyer/products/{spu_no}" method:"get" tags:"Catalog-Buyer" summary:"Get product detail"`
	SpuNo  string `json:"spu_no"`
}
type GetProductDetailRes = pb.GetProductDetailRes

type ListProductsReq struct {
	g.Meta     `path:"/v1/catalog/buyer/products" method:"get" tags:"Catalog-Buyer" summary:"List products"`
	CategoryId uint64    `json:"category_id"`
	Page       int32     `json:"page"`
	PageSize   int32     `json:"page_size"`
	SortBy     pb.SortBy `json:"sort_by"`
}
type ListProductsRes = pb.ListProductsRes

type SearchProductsReq struct {
	g.Meta     `path:"/v1/catalog/buyer/products/search" method:"get" tags:"Catalog-Buyer" summary:"Search products"`
	Keyword    string    `json:"keyword"`
	CategoryId uint64    `json:"category_id"`
	Page       int32     `json:"page"`
	PageSize   int32     `json:"page_size"`
	SortBy     pb.SortBy `json:"sort_by"`
}
type SearchProductsRes = pb.SearchProductsRes

type ListBuyerProductImagesReq struct {
	g.Meta `path:"/v1/catalog/buyer/product-images" method:"get" tags:"Catalog-Buyer" summary:"List product images"`
	SpuNos string `json:"spu_nos" dc:"comma separated spu_no list"`
}

type BuyerProductImageItem struct {
	SpuNo    string `json:"spu_no"`
	ImageUrl string `json:"image_url"`
}

type ListBuyerProductImagesRes struct {
	Items []BuyerProductImageItem `json:"items"`
}
