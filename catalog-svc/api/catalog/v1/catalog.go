package v1

import (
	pb "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateProductDraftReq struct {
	g.Meta              `path:"/v1/catalog/seller/products/draft" method:"post" tags:"Catalog-Seller" summary:"Create product draft"`
	ShopNo              string                 `json:"shop_no" p:"shop_no"`
	Title               string                 `json:"title" p:"title"`
	SubTitle            string                 `json:"sub_title" p:"sub_title"`
	CategoryId          uint64                 `json:"category_id" p:"category_id"`
	BrandNo             string                 `json:"brand_no" p:"brand_no"`
	MainImageAssetIds   []uint64               `json:"main_image_asset_ids" p:"main_image_asset_ids"`
	DetailImageAssetIds []uint64               `json:"detail_image_asset_ids" p:"detail_image_asset_ids"`
	SpuAttrs            []*pb.AttributeValue   `json:"spu_attrs" p:"spu_attrs"`
	PublishTime         *timestamppb.Timestamp `json:"publish_time" p:"publish_time"`
}
type CreateProductDraftRes = pb.CreateProductDraftRes

type UpdateProductDraftReq struct {
	g.Meta          `path:"/v1/catalog/seller/products/draft" method:"put" tags:"Catalog-Seller" summary:"Update product draft"`
	SpuNo           string                 `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32                 `json:"expected_version" p:"expected_version"`
	Patch           *pb.ProductDraftPatch  `json:"patch" p:"patch"`
	UpdateMask      *fieldmaskpb.FieldMask `json:"update_mask" p:"update_mask"`
}
type UpdateProductDraftRes = pb.UpdateProductDraftRes

type UpsertSkuDraftsReq struct {
	g.Meta          `path:"/v1/catalog/seller/products/skus:upsert" method:"post" tags:"Catalog-Seller" summary:"Upsert sku drafts"`
	SpuNo           string             `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32             `json:"expected_version" p:"expected_version"`
	Items           []*pb.SkuDraftItem `json:"items" p:"items"`
	ReplaceAll      bool               `json:"replace_all" p:"replace_all"`
}
type UpsertSkuDraftsRes = pb.UpsertSkuDraftsRes

type SubmitProductReviewReq struct {
	g.Meta          `path:"/v1/catalog/seller/products/review:submit" method:"post" tags:"Catalog-Seller" summary:"Submit product review"`
	SpuNo           string `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32 `json:"expected_version" p:"expected_version"`
	SubmitNote      string `json:"submit_note" p:"submit_note"`
}
type SubmitProductReviewRes = pb.SubmitProductReviewRes

type ResubmitProductReviewReq struct {
	g.Meta          `path:"/v1/catalog/seller/products/review:resubmit" method:"post" tags:"Catalog-Seller" summary:"Resubmit product review"`
	SpuNo           string `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32 `json:"expected_version" p:"expected_version"`
	SubmitNote      string `json:"submit_note" p:"submit_note"`
}
type ResubmitProductReviewRes = pb.ResubmitProductReviewRes

type SetProductOnShelfReq struct {
	g.Meta               `path:"/v1/catalog/seller/products/on-shelf" method:"post" tags:"Catalog-Seller" summary:"Set product on shelf"`
	SpuNo                string                 `json:"spu_no" p:"spu_no"`
	ExpectedVersion      uint32                 `json:"expected_version" p:"expected_version"`
	EffectivePublishTime *timestamppb.Timestamp `json:"effective_publish_time" p:"effective_publish_time"`
}
type SetProductOnShelfRes = pb.SetProductOnShelfRes

type SetProductOffShelfReq struct {
	g.Meta          `path:"/v1/catalog/seller/products/off-shelf" method:"post" tags:"Catalog-Seller" summary:"Set product off shelf"`
	SpuNo           string `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32 `json:"expected_version" p:"expected_version"`
	ReasonCode      string `json:"reason_code" p:"reason_code"`
}
type SetProductOffShelfRes = pb.SetProductOffShelfRes

type DeleteProductDraftReq struct {
	g.Meta          `path:"/v1/catalog/seller/products/draft:delete" method:"post" tags:"Catalog-Seller" summary:"Delete product draft"`
	SpuNo           string `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32 `json:"expected_version" p:"expected_version"`
	ReasonCode      string `json:"reason_code" p:"reason_code"`
}
type DeleteProductDraftRes struct {
	Success bool `json:"success"`
}

type GetMyProductReq struct {
	g.Meta `path:"/v1/catalog/seller/products/{spu_no}" method:"get" tags:"Catalog-Seller" summary:"Get my product detail"`
	SpuNo  string `json:"spu_no" p:"spu_no" in:"path" dc:"spu no"`
}
type GetMyProductRes = pb.GetMyProductRes

type ListMyProductsReq struct {
	g.Meta   `path:"/v1/catalog/seller/products" method:"get" tags:"Catalog-Seller" summary:"List my products"`
	Page     int32          `json:"page" p:"page" in:"query"`
	PageSize int32          `json:"page_size" p:"page_size" in:"query"`
	Statuses []pb.SpuStatus `json:"statuses" p:"statuses" in:"query"`
	Keyword  string         `json:"keyword" p:"keyword" in:"query"`
}
type ListMyProductsRes = pb.ListMyProductsRes

type ListReviewTasksReq struct {
	g.Meta   `path:"/v1/catalog/admin/review/tasks" method:"get" tags:"Catalog-Admin" summary:"List review tasks"`
	Page     int32          `json:"page" p:"page" in:"query"`
	PageSize int32          `json:"page_size" p:"page_size" in:"query"`
	Statuses []pb.SpuStatus `json:"statuses" p:"statuses" in:"query"`
	Keyword  string         `json:"keyword" p:"keyword" in:"query"`
}
type ListReviewTasksRes = pb.ListReviewTasksRes

type GetReviewDetailReq struct {
	g.Meta `path:"/v1/catalog/admin/review/{spu_no}" method:"get" tags:"Catalog-Admin" summary:"Get review detail"`
	SpuNo  string `json:"spu_no" p:"spu_no" in:"path"`
}
type GetReviewDetailRes = pb.GetReviewDetailRes

type ApproveProductReq struct {
	g.Meta          `path:"/v1/catalog/admin/review/approve" method:"post" tags:"Catalog-Admin" summary:"Approve product"`
	SpuNo           string `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32 `json:"expected_version" p:"expected_version"`
	ReviewComment   string `json:"review_comment" p:"review_comment"`
}
type ApproveProductRes = pb.ApproveProductRes

type RejectProductReq struct {
	g.Meta           `path:"/v1/catalog/admin/review/reject" method:"post" tags:"Catalog-Admin" summary:"Reject product"`
	SpuNo            string `json:"spu_no" p:"spu_no"`
	ExpectedVersion  uint32 `json:"expected_version" p:"expected_version"`
	RejectReasonCode string `json:"reject_reason_code" p:"reject_reason_code"`
	RejectComment    string `json:"reject_comment" p:"reject_comment"`
}
type RejectProductRes = pb.RejectProductRes

type FreezeProductReq struct {
	g.Meta          `path:"/v1/catalog/admin/review/freeze" method:"post" tags:"Catalog-Admin" summary:"Freeze product"`
	SpuNo           string `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32 `json:"expected_version" p:"expected_version"`
	ReasonCode      string `json:"reason_code" p:"reason_code"`
	Reason          string `json:"reason" p:"reason"`
}
type FreezeProductRes = pb.FreezeProductRes

type UnfreezeProductReq struct {
	g.Meta          `path:"/v1/catalog/admin/review/unfreeze" method:"post" tags:"Catalog-Admin" summary:"Unfreeze product"`
	SpuNo           string `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32 `json:"expected_version" p:"expected_version"`
}
type UnfreezeProductRes = pb.UnfreezeProductRes

type ForceOffShelfReq struct {
	g.Meta          `path:"/v1/catalog/admin/review/force-off-shelf" method:"post" tags:"Catalog-Admin" summary:"Force off shelf"`
	SpuNo           string `json:"spu_no" p:"spu_no"`
	ExpectedVersion uint32 `json:"expected_version" p:"expected_version"`
	ReasonCode      string `json:"reason_code" p:"reason_code"`
	Reason          string `json:"reason" p:"reason"`
}
type ForceOffShelfRes = pb.ForceOffShelfRes

type GetProductDetailReq struct {
	g.Meta `path:"/v1/catalog/buyer/products/{spu_no}" method:"get" tags:"Catalog-Buyer" summary:"Get product detail"`
	SpuNo  string `json:"spu_no" p:"spu_no" in:"path"`
}
type GetProductDetailRes = pb.GetProductDetailRes

type ListProductsReq struct {
	g.Meta     `path:"/v1/catalog/buyer/products" method:"get" tags:"Catalog-Buyer" summary:"List products"`
	CategoryId uint64    `json:"category_id" p:"category_id" in:"query"`
	Page       int32     `json:"page" p:"page" in:"query"`
	PageSize   int32     `json:"page_size" p:"page_size" in:"query"`
	SortBy     pb.SortBy `json:"sort_by" p:"sort_by" in:"query"`
}
type ListProductsRes = pb.ListProductsRes

type SearchProductsReq struct {
	g.Meta     `path:"/v1/catalog/buyer/products/search" method:"get" tags:"Catalog-Buyer" summary:"Search products"`
	Keyword    string    `json:"keyword" p:"keyword" in:"query"`
	CategoryId uint64    `json:"category_id" p:"category_id" in:"query"`
	Page       int32     `json:"page" p:"page" in:"query"`
	PageSize   int32     `json:"page_size" p:"page_size" in:"query"`
	SortBy     pb.SortBy `json:"sort_by" p:"sort_by" in:"query"`
}
type SearchProductsRes = pb.SearchProductsRes

type ListBuyerProductImagesReq struct {
	g.Meta `path:"/v1/catalog/buyer/product-images" method:"get" tags:"Catalog-Buyer" summary:"List product images"`
	SpuNos string `json:"spu_nos" p:"spu_nos" in:"query" dc:"comma separated spu_no list"`
}

type BuyerProductImageItem struct {
	SpuNo    string `json:"spu_no"`
	ImageUrl string `json:"image_url"`
}

type ListBuyerProductImagesRes struct {
	Items []BuyerProductImageItem `json:"items"`
}

type UpsertSkuStockProjectionInternalReq struct {
	g.Meta `path:"/v1/catalog/internal/stock-projection:upsert" method:"post" tags:"Catalog-Internal" summary:"Upsert sku stock projection"`
	Items  []*pb.UpsertSkuStockProjectionItem `json:"items" p:"items"`
}

type UpsertSkuStockProjectionInternalRes = pb.UpsertSkuStockProjectionRes
