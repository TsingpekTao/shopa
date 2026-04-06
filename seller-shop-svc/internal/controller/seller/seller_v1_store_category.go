package seller

import (
	"context"

	sellerv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/seller/v1"
	pb "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
)

func (c *ControllerV1) ListStoreCategories(ctx context.Context, req *sellerv1.ListStoreCategoriesReq) (*sellerv1.ListStoreCategoriesRes, error) {
	out, err := c.svc.ListStoreCategories(withRequestMetadata(ctx), &pb.ListStoreCategoriesReq{ShopNo: req.ShopNo, BuyerSide: req.BuyerSide})
	if err != nil {
		return nil, err
	}
	return &sellerv1.ListStoreCategoriesRes{Categories: out.Categories}, nil
}

func (c *ControllerV1) CreateStoreCategory(ctx context.Context, req *sellerv1.CreateStoreCategoryReq) (*sellerv1.CreateStoreCategoryRes, error) {
	isVisible := true
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}
	out, err := c.svc.CreateStoreCategory(withRequestMetadata(ctx), &pb.CreateStoreCategoryReq{
		ShopNo:    req.ShopNo,
		ParentId:  req.ParentId,
		Name:      req.Name,
		SortOrder: req.SortOrder,
		IsVisible: isVisible,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.CreateStoreCategoryRes{Category: out.Category}, nil
}

func (c *ControllerV1) UpdateStoreCategory(ctx context.Context, req *sellerv1.UpdateStoreCategoryReq) (*sellerv1.UpdateStoreCategoryRes, error) {
	var sortOrder int32
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	isVisible := true
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}
	out, err := c.svc.UpdateStoreCategory(withRequestMetadata(ctx), &pb.UpdateStoreCategoryReq{
		ShopNo:       req.ShopNo,
		CategoryId:   req.CategoryId,
		Name:         req.Name,
		SortOrder:    sortOrder,
		IsVisible:    isVisible,
		SetSortOrder: req.SortOrder != nil,
		SetIsVisible: req.IsVisible != nil,
	})
	if err != nil {
		return nil, err
	}
	return &sellerv1.UpdateStoreCategoryRes{Category: out.Category}, nil
}

func (c *ControllerV1) SortStoreCategories(ctx context.Context, req *sellerv1.SortStoreCategoriesReq) (*sellerv1.SortStoreCategoriesRes, error) {
	out, err := c.svc.SortStoreCategories(withRequestMetadata(ctx), &pb.SortStoreCategoriesReq{ShopNo: req.ShopNo, Items: req.Items})
	if err != nil {
		return nil, err
	}
	return &sellerv1.SortStoreCategoriesRes{Categories: out.Categories}, nil
}

func (c *ControllerV1) DeleteStoreCategory(ctx context.Context, req *sellerv1.DeleteStoreCategoryReq) (*sellerv1.DeleteStoreCategoryRes, error) {
	out, err := c.svc.DeleteStoreCategory(withRequestMetadata(ctx), &pb.DeleteStoreCategoryReq{ShopNo: req.ShopNo, CategoryId: req.CategoryId})
	if err != nil {
		return nil, err
	}
	return &sellerv1.DeleteStoreCategoryRes{Ok: out.Ok}, nil
}

func (c *ControllerV1) GetProductStoreCategoryBinding(ctx context.Context, req *sellerv1.GetProductStoreCategoryBindingReq) (*sellerv1.GetProductStoreCategoryBindingRes, error) {
	out, err := c.svc.GetProductStoreCategoryBinding(withRequestMetadata(ctx), &pb.GetProductStoreCategoryBindingReq{ShopNo: req.ShopNo, SpuNo: req.SpuNo})
	if err != nil {
		return nil, err
	}
	return &sellerv1.GetProductStoreCategoryBindingRes{Binding: out.Binding}, nil
}

func (c *ControllerV1) BatchGetProductStoreCategoryBindings(ctx context.Context, req *sellerv1.BatchGetProductStoreCategoryBindingsReq) (*sellerv1.BatchGetProductStoreCategoryBindingsRes, error) {
	out, err := c.svc.BatchGetProductStoreCategoryBindings(withRequestMetadata(ctx), &pb.BatchGetProductStoreCategoryBindingsReq{ShopNo: req.ShopNo, SpuNos: req.SpuNos})
	if err != nil {
		return nil, err
	}
	return &sellerv1.BatchGetProductStoreCategoryBindingsRes{Items: out.Items}, nil
}

func (c *ControllerV1) UpdateProductStoreCategoryBinding(ctx context.Context, req *sellerv1.UpdateProductStoreCategoryBindingReq) (*sellerv1.UpdateProductStoreCategoryBindingRes, error) {
	out, err := c.svc.UpdateProductStoreCategoryBinding(withRequestMetadata(ctx), &pb.UpdateProductStoreCategoryBindingReq{ShopNo: req.ShopNo, SpuNo: req.SpuNo, StoreCategoryId: req.StoreCategoryId})
	if err != nil {
		return nil, err
	}
	return &sellerv1.UpdateProductStoreCategoryBindingRes{Binding: out.Binding}, nil
}

func (c *ControllerV1) ListBuyerStoreCategories(ctx context.Context, req *sellerv1.ListBuyerStoreCategoriesReq) (*sellerv1.ListBuyerStoreCategoriesRes, error) {
	out, err := c.svc.ListBuyerStoreCategories(withRequestMetadata(ctx), &pb.ListBuyerStoreCategoriesReq{ShopNo: req.ShopNo})
	if err != nil {
		return nil, err
	}
	return &sellerv1.ListBuyerStoreCategoriesRes{Categories: out.Categories}, nil
}
