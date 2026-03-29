export type BuyerSortBy = 0 | 1 | 2 | 3 | 4 | 5;

export const BUYER_SORT_BY = {
  UNSPECIFIED: 0 as BuyerSortBy,
  DEFAULT: 1 as BuyerSortBy,
  PRICE_ASC: 2 as BuyerSortBy,
  PRICE_DESC: 3 as BuyerSortBy,
  SALES_DESC: 4 as BuyerSortBy,
  NEWEST: 5 as BuyerSortBy
};

export type BuyerProductCard = {
  spuNo: string;
  title: string;
  categoryId: number;
  brandNo: string;
  coverImageAssetId: string;
  minSalePrice: number;
  maxSalePrice: number;
  soldCount: number;
};

export type ListBuyerProductsParams = {
  categoryId?: number;
  page?: number;
  pageSize?: number;
  sortBy?: BuyerSortBy;
};

export type ListBuyerProductsResult = {
  items: BuyerProductCard[];
  page: number;
  pageSize: number;
  total: number;
};

export type BuyerProductImageItem = {
  spuNo: string;
  imageUrl: string;
};
