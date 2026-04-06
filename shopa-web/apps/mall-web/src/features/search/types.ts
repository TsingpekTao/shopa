export type SearchProductCard = {
  spuNo: string;
  title: string;
  coverUrl: string;
  minPrice: number;
  maxPrice: number;
  categoryId?: number;
  categoryName?: string;
  shopNo: string;
  shopName: string;
  salesCount: number;
  stockTotal: number;
  avgScoreX100: number;
  reviewTotal: number;
};

export type SearchShopPreviewItem = {
  spuNo: string;
  title: string;
  coverUrl: string;
  minPrice: number;
  salesCount: number;
};

export type SearchShopCard = {
  shopNo: string;
  shopName: string;
  shopDisplayName: string;
  coverUrl: string;
  matchedProductCount: number;
  sampleSpuNo: string;
  previewItems: SearchShopPreviewItem[];
};

export type SearchProductsResult = {
  items: SearchProductCard[];
  shops: SearchShopCard[];
  nextCursor: string;
  hasMore: boolean;
  source: "search-svc" | "catalog-fallback";
};
