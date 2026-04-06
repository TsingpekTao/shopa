package sellershopv1

func (x *ListStoreCategoriesReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}

func (x *ListStoreCategoriesReq) GetBuyerSide() bool {
	if x != nil {
		return x.BuyerSide
	}
	return false
}

func (x *CreateStoreCategoryReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}
func (x *CreateStoreCategoryReq) GetParentId() uint64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}
func (x *CreateStoreCategoryReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}
func (x *CreateStoreCategoryReq) GetSortOrder() int32 {
	if x != nil {
		return x.SortOrder
	}
	return 0
}
func (x *CreateStoreCategoryReq) GetIsVisible() bool {
	if x != nil {
		return x.IsVisible
	}
	return false
}

func (x *UpdateStoreCategoryReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}
func (x *UpdateStoreCategoryReq) GetCategoryId() uint64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}
func (x *UpdateStoreCategoryReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}
func (x *UpdateStoreCategoryReq) GetSortOrder() int32 {
	if x != nil {
		return x.SortOrder
	}
	return 0
}
func (x *UpdateStoreCategoryReq) GetIsVisible() bool {
	if x != nil {
		return x.IsVisible
	}
	return false
}

func (x *StoreCategorySortItem) GetCategoryId() uint64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}
func (x *StoreCategorySortItem) GetSortOrder() int32 {
	if x != nil {
		return x.SortOrder
	}
	return 0
}

func (x *SortStoreCategoriesReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}
func (x *SortStoreCategoriesReq) GetItems() []*StoreCategorySortItem {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *DeleteStoreCategoryReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}
func (x *DeleteStoreCategoryReq) GetCategoryId() uint64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

func (x *GetProductStoreCategoryBindingReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}
func (x *GetProductStoreCategoryBindingReq) GetSpuNo() string {
	if x != nil {
		return x.SpuNo
	}
	return ""
}

func (x *BatchGetProductStoreCategoryBindingsReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}
func (x *BatchGetProductStoreCategoryBindingsReq) GetSpuNos() []string {
	if x != nil {
		return x.SpuNos
	}
	return nil
}

func (x *UpdateProductStoreCategoryBindingReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}
func (x *UpdateProductStoreCategoryBindingReq) GetSpuNo() string {
	if x != nil {
		return x.SpuNo
	}
	return ""
}
func (x *UpdateProductStoreCategoryBindingReq) GetStoreCategoryId() uint64 {
	if x != nil {
		return x.StoreCategoryId
	}
	return 0
}

func (x *ListBuyerStoreCategoriesReq) GetShopNo() string {
	if x != nil {
		return x.ShopNo
	}
	return ""
}
