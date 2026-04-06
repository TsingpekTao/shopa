package sellershopv1

import timestamppb "google.golang.org/protobuf/types/known/timestamppb"

type StoreCategory struct {
	Id           uint64                 `json:"id,omitempty"`
	ShopNo       string                 `json:"shop_no,omitempty"`
	ParentId     uint64                 `json:"parent_id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Level        int32                  `json:"level,omitempty"`
	SortOrder    int32                  `json:"sort_order,omitempty"`
	IsVisible    bool                   `json:"is_visible,omitempty"`
	IsDeleted    bool                   `json:"is_deleted,omitempty"`
	ProductCount int32                  `json:"product_count,omitempty"`
	Children     []*StoreCategory       `json:"children,omitempty"`
	CreatedAt    *timestamppb.Timestamp `json:"created_at,omitempty"`
	UpdatedAt    *timestamppb.Timestamp `json:"updated_at,omitempty"`
}

type ProductStoreCategoryBinding struct {
	ShopNo            string                 `json:"shop_no,omitempty"`
	SpuNo             string                 `json:"spu_no,omitempty"`
	StoreCategoryId   uint64                 `json:"store_category_id,omitempty"`
	StoreCategoryL1   uint64                 `json:"store_category_l1,omitempty"`
	StoreCategoryL2   uint64                 `json:"store_category_l2,omitempty"`
	StoreCategoryPath []uint64               `json:"store_category_path,omitempty"`
	StoreCategoryName string                 `json:"store_category_name,omitempty"`
	UpdatedAt         *timestamppb.Timestamp `json:"updated_at,omitempty"`
}

type ListStoreCategoriesReq struct {
	ShopNo    string `json:"shop_no,omitempty"`
	BuyerSide bool   `json:"buyer_side,omitempty"`
}

type ListStoreCategoriesRes struct {
	Categories []*StoreCategory `json:"categories,omitempty"`
}

type ListBuyerStoreCategoriesReq struct {
	ShopNo string `json:"shop_no,omitempty"`
}

type ListBuyerStoreCategoriesRes struct {
	Categories []*StoreCategory `json:"categories,omitempty"`
}

type CreateStoreCategoryReq struct {
	ShopNo    string `json:"shop_no,omitempty"`
	ParentId  uint64 `json:"parent_id,omitempty"`
	Name      string `json:"name,omitempty"`
	SortOrder int32  `json:"sort_order,omitempty"`
	IsVisible bool   `json:"is_visible,omitempty"`
}

type CreateStoreCategoryRes struct {
	Category *StoreCategory `json:"category,omitempty"`
}

type UpdateStoreCategoryReq struct {
	ShopNo       string `json:"shop_no,omitempty"`
	CategoryId   uint64 `json:"category_id,omitempty"`
	Name         string `json:"name,omitempty"`
	SortOrder    int32  `json:"sort_order,omitempty"`
	IsVisible    bool   `json:"is_visible,omitempty"`
	SetSortOrder bool   `json:"-"`
	SetIsVisible bool   `json:"-"`
}

type UpdateStoreCategoryRes struct {
	Category *StoreCategory `json:"category,omitempty"`
}

type StoreCategorySortItem struct {
	CategoryId uint64 `json:"category_id,omitempty"`
	SortOrder  int32  `json:"sort_order,omitempty"`
}

type SortStoreCategoriesReq struct {
	ShopNo string                   `json:"shop_no,omitempty"`
	Items  []*StoreCategorySortItem `json:"items,omitempty"`
}

type SortStoreCategoriesRes struct {
	Categories []*StoreCategory `json:"categories,omitempty"`
}

type DeleteStoreCategoryReq struct {
	ShopNo     string `json:"shop_no,omitempty"`
	CategoryId uint64 `json:"category_id,omitempty"`
}

type DeleteStoreCategoryRes struct {
	Ok bool `json:"ok,omitempty"`
}

type GetProductStoreCategoryBindingReq struct {
	ShopNo string `json:"shop_no,omitempty"`
	SpuNo  string `json:"spu_no,omitempty"`
}

type GetProductStoreCategoryBindingRes struct {
	Binding *ProductStoreCategoryBinding `json:"binding,omitempty"`
}

type BatchGetProductStoreCategoryBindingsReq struct {
	ShopNo string   `json:"shop_no,omitempty"`
	SpuNos []string `json:"spu_nos,omitempty"`
}

type BatchGetProductStoreCategoryBindingsRes struct {
	Items []*ProductStoreCategoryBinding `json:"items,omitempty"`
}

type UpdateProductStoreCategoryBindingReq struct {
	ShopNo          string `json:"shop_no,omitempty"`
	SpuNo           string `json:"spu_no,omitempty"`
	StoreCategoryId uint64 `json:"store_category_id,omitempty"`
}

type UpdateProductStoreCategoryBindingRes struct {
	Binding *ProductStoreCategoryBinding `json:"binding,omitempty"`
}

type ReconcileStoreCategoryCountsReq struct{}

type ReconcileStoreCategoryCountsRes struct {
	UpdatedCategories int32 `json:"updated_categories,omitempty"`
}
