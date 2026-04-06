export type MallShopDetail = {
  shopNo: string;
  shopName: string;
  shopDisplayName: string;
  status: number;
  buyerVisible: boolean;
};

export type ShopStoreCategory = {
  id: number;
  shopNo: string;
  parentId: number;
  name: string;
  level: number;
  sortOrder: number;
  productCount: number;
  isVisible: boolean;
  children: ShopStoreCategory[];
};
