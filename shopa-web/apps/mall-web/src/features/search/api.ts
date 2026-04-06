import { SearchProductCard, SearchProductsResult, SearchShopCard, SearchShopPreviewItem } from "./types";

type RawSearchItem = {
  spuNo?: string;
  spu_no?: string;
  title?: string;
  coverUrl?: string;
  cover_url?: string;
  minPrice?: number | string;
  min_price?: number | string;
  maxPrice?: number | string;
  max_price?: number | string;
  categoryId?: number | string;
  category_id?: number | string;
  categoryName?: string;
  category_name?: string;
  shopNo?: string;
  shop_no?: string;
  shopName?: string;
  shop_name?: string;
  salesCount?: number | string;
  sales_count?: number | string;
  stockTotal?: number | string;
  stock_total?: number | string;
  avgScoreX100?: number | string;
  avg_score_x100?: number | string;
  reviewTotal?: number | string;
  review_total?: number | string;
};

type RawSearchShop = {
  shopNo?: string;
  shop_no?: string;
  shopName?: string;
  shop_name?: string;
  shopDisplayName?: string;
  shop_display_name?: string;
  coverUrl?: string;
  cover_url?: string;
  matchedProductCount?: number | string;
  matched_product_count?: number | string;
  sampleSpuNo?: string;
  sample_spu_no?: string;
  previewItems?: RawSearchShopPreviewItem[];
  preview_items?: RawSearchShopPreviewItem[];
};

type RawSearchShopPreviewItem = {
  spuNo?: string;
  spu_no?: string;
  title?: string;
  coverUrl?: string;
  cover_url?: string;
  minPrice?: number | string;
  min_price?: number | string;
  salesCount?: number | string;
  sales_count?: number | string;
};

type RawSearchProductsResult = {
  items?: RawSearchItem[];
  shops?: RawSearchShop[];
  nextCursor?: string;
  next_cursor?: string;
  hasMore?: boolean;
  has_more?: boolean;
  source?: "search-svc" | "catalog-fallback";
};

function toNumber(value: unknown, fallback = 0): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function normalizeItem(raw: RawSearchItem): SearchProductCard {
  return {
    spuNo: toString(raw.spuNo ?? raw.spu_no),
    title: toString(raw.title),
    coverUrl: toString(raw.coverUrl ?? raw.cover_url),
    minPrice: toNumber(raw.minPrice ?? raw.min_price),
    maxPrice: toNumber(raw.maxPrice ?? raw.max_price),
    categoryId: toNumber(raw.categoryId ?? raw.category_id, 0) || undefined,
    categoryName: toString(raw.categoryName ?? raw.category_name) || undefined,
    shopNo: toString(raw.shopNo ?? raw.shop_no),
    shopName: toString(raw.shopName ?? raw.shop_name),
    salesCount: toNumber(raw.salesCount ?? raw.sales_count),
    stockTotal: toNumber(raw.stockTotal ?? raw.stock_total),
    avgScoreX100: toNumber(raw.avgScoreX100 ?? raw.avg_score_x100),
    reviewTotal: toNumber(raw.reviewTotal ?? raw.review_total)
  };
}

function normalizeShopPreviewItem(raw: RawSearchShopPreviewItem): SearchShopPreviewItem {
  return {
    spuNo: toString(raw.spuNo ?? raw.spu_no),
    title: toString(raw.title),
    coverUrl: toString(raw.coverUrl ?? raw.cover_url),
    minPrice: toNumber(raw.minPrice ?? raw.min_price),
    salesCount: toNumber(raw.salesCount ?? raw.sales_count)
  };
}

function normalizeShop(raw: RawSearchShop): SearchShopCard {
  return {
    shopNo: toString(raw.shopNo ?? raw.shop_no),
    shopName: toString(raw.shopName ?? raw.shop_name),
    shopDisplayName: toString(raw.shopDisplayName ?? raw.shop_display_name),
    coverUrl: toString(raw.coverUrl ?? raw.cover_url),
    matchedProductCount: toNumber(raw.matchedProductCount ?? raw.matched_product_count),
    sampleSpuNo: toString(raw.sampleSpuNo ?? raw.sample_spu_no),
    previewItems: Array.isArray(raw.previewItems ?? raw.preview_items)
      ? (raw.previewItems ?? raw.preview_items ?? []).map((item) => normalizeShopPreviewItem(item))
      : []
  };
}

export async function searchProducts(params: {
  query: string;
  pageSize?: number;
  nextCursor?: string;
  shopNo?: string;
  storeCategoryId?: number;
  storeCategoryLevel?: 1 | 2;
}): Promise<SearchProductsResult> {
  const searchParams = new URLSearchParams();
  if (params.query.trim()) {
    searchParams.set("q", params.query.trim());
  }
  if (params.pageSize && params.pageSize > 0) {
    searchParams.set("pageSize", String(params.pageSize));
  }
  if (params.nextCursor?.trim()) {
    searchParams.set("nextCursor", params.nextCursor.trim());
  }
  if (params.shopNo?.trim()) {
    searchParams.set("shopNo", params.shopNo.trim());
  }
  if (params.storeCategoryId && params.storeCategoryId > 0) {
    searchParams.set("storeCategoryId", String(params.storeCategoryId));
    if (params.storeCategoryLevel) {
      searchParams.set("storeCategoryLevel", String(params.storeCategoryLevel));
    }
  }

  const response = await fetch(`/api/search/products?${searchParams.toString()}`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    throw new Error(`search request failed: ${response.status}`);
  }

  const payload = (await response.json()) as RawSearchProductsResult;
  return {
    items: Array.isArray(payload.items) ? payload.items.map((item) => normalizeItem(item)) : [],
    shops: Array.isArray(payload.shops) ? payload.shops.map((item) => normalizeShop(item)) : [],
    nextCursor: toString(payload.nextCursor ?? payload.next_cursor),
    hasMore: Boolean(payload.hasMore ?? payload.has_more),
    source: payload.source === "catalog-fallback" ? "catalog-fallback" : "search-svc"
  };
}
