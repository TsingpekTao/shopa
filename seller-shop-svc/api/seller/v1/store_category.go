package v1

import (
	pb "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type ListStoreCategoriesReq struct {
	g.Meta    `path:"/v1/seller/shops/{shopNo}/store-categories" method:"get" tags:"SellerShop" summary:"List store categories for a shop"`
	ShopNo    string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	BuyerSide bool   `json:"buyerSide" in:"query"`
}

type ListStoreCategoriesRes struct {
	Categories []*pb.StoreCategory `json:"categories"`
}

type CreateStoreCategoryReq struct {
	g.Meta    `path:"/v1/seller/shops/{shopNo}/store-categories" method:"post" tags:"SellerShop" summary:"Create a store category"`
	ShopNo    string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	ParentId  uint64 `json:"parentId"`
	Name      string `json:"name" v:"required#name is required"`
	SortOrder int32  `json:"sortOrder"`
	IsVisible *bool  `json:"isVisible"`
}

type CreateStoreCategoryRes struct {
	Category *pb.StoreCategory `json:"category"`
}

type UpdateStoreCategoryReq struct {
	g.Meta     `path:"/v1/seller/shops/{shopNo}/store-categories/{categoryId}" method:"patch" tags:"SellerShop" summary:"Update a store category"`
	ShopNo     string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	CategoryId uint64 `json:"categoryId" in:"path" v:"required#categoryId is required"`
	Name       string `json:"name"`
	SortOrder  *int32 `json:"sortOrder"`
	IsVisible  *bool  `json:"isVisible"`
}

type UpdateStoreCategoryRes struct {
	Category *pb.StoreCategory `json:"category"`
}

type SortStoreCategoriesReq struct {
	g.Meta `path:"/v1/seller/shops/{shopNo}/store-categories:sort" method:"post" tags:"SellerShop" summary:"Sort store categories"`
	ShopNo string                      `json:"shopNo" in:"path" v:"required#shopNo is required"`
	Items  []*pb.StoreCategorySortItem `json:"items" v:"required#items is required"`
}

type SortStoreCategoriesRes struct {
	Categories []*pb.StoreCategory `json:"categories"`
}

type DeleteStoreCategoryReq struct {
	g.Meta     `path:"/v1/seller/shops/{shopNo}/store-categories/{categoryId}" method:"delete" tags:"SellerShop" summary:"Soft delete an empty store category"`
	ShopNo     string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	CategoryId uint64 `json:"categoryId" in:"path" v:"required#categoryId is required"`
}

type DeleteStoreCategoryRes struct {
	Ok bool `json:"ok"`
}

type GetProductStoreCategoryBindingReq struct {
	g.Meta `path:"/v1/seller/shops/{shopNo}/products/{spuNo}/store-category" method:"get" tags:"SellerShop" summary:"Get product store category binding"`
	ShopNo string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	SpuNo  string `json:"spuNo" in:"path" v:"required#spuNo is required"`
}

type GetProductStoreCategoryBindingRes struct {
	Binding *pb.ProductStoreCategoryBinding `json:"binding"`
}

type BatchGetProductStoreCategoryBindingsReq struct {
	g.Meta `path:"/v1/seller/shops/{shopNo}/products/store-categories:batch-get" method:"post" tags:"SellerShop" summary:"Batch get product store category bindings"`
	ShopNo string   `json:"shopNo" in:"path" v:"required#shopNo is required"`
	SpuNos []string `json:"spuNos" v:"required#spuNos is required"`
}

type BatchGetProductStoreCategoryBindingsRes struct {
	Items []*pb.ProductStoreCategoryBinding `json:"items"`
}

type UpdateProductStoreCategoryBindingReq struct {
	g.Meta          `path:"/v1/seller/shops/{shopNo}/products/{spuNo}/store-category" method:"put" tags:"SellerShop" summary:"Update product store category binding"`
	ShopNo          string `json:"shopNo" in:"path" v:"required#shopNo is required"`
	SpuNo           string `json:"spuNo" in:"path" v:"required#spuNo is required"`
	StoreCategoryId uint64 `json:"storeCategoryId"`
}

type UpdateProductStoreCategoryBindingRes struct {
	Binding *pb.ProductStoreCategoryBinding `json:"binding"`
}

type ListBuyerStoreCategoriesReq struct {
	g.Meta `path:"/v1/shops/{shopNo}/store-categories" method:"get" tags:"SellerShopBuyer" summary:"List buyer-visible store categories"`
	ShopNo string `json:"shopNo" in:"path" v:"required#shopNo is required"`
}

type ListBuyerStoreCategoriesRes struct {
	Categories []*pb.StoreCategory `json:"categories"`
}
