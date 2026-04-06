import { NextRequest, NextResponse } from "next/server";
import { resolveBuyerCategoryDescriptor } from "@/features/search/category";

type SearchSvcItem = {
  spu_no?: string;
  spuNo?: string;
  title?: string;
  cover_asset_id?: number | string;
  coverAssetId?: number | string;
  cover_url?: string;
  coverUrl?: string;
  min_price?: number | string;
  minPrice?: number | string;
  max_price?: number | string;
  maxPrice?: number | string;
  category_id?: number | string;
  categoryId?: number | string;
  category_name?: string;
  categoryName?: string;
  shop_no?: string;
  shopNo?: string;
  shop_name?: string;
  shopName?: string;
  sales_count?: number | string;
  salesCount?: number | string;
  stock_total?: number | string;
  stockTotal?: number | string;
  avg_score_x100?: number | string;
  avgScoreX100?: number | string;
  review_total?: number | string;
  reviewTotal?: number | string;
};

type SearchSvcResponse = {
  list?: SearchSvcItem[];
  next_cursor?: string;
  nextCursor?: string;
  has_more?: boolean;
  hasMore?: boolean;
};

type CatalogItem = {
  spu_no?: string;
  spuNo?: string;
  title?: string;
  sub_title?: string;
  subTitle?: string;
  min_sale_price?: number | string;
  minSalePrice?: number | string;
  max_sale_price?: number | string;
  maxSalePrice?: number | string;
  sold_count?: number | string;
  soldCount?: number | string;
  cover_image_asset_id?: number | string;
  coverImageAssetId?: number | string;
  main_image_asset_ids?: Array<number | string>;
  mainImageAssetIds?: Array<number | string>;
  category_id?: number | string;
  categoryId?: number | string;
  category_name?: string;
  categoryName?: string;
  shop_no?: string;
  shopNo?: string;
  spu_status?: number | string;
  spuStatus?: number | string;
};

type CatalogSearchResponse = {
  items?: CatalogItem[];
  products?: CatalogItem[];
  page?: number;
  page_size?: number;
  pageSize?: number;
  total?: number;
};

type ProductImageItem = {
  spu_no?: string;
  spuNo?: string;
  image_url?: string;
  imageUrl?: string;
};

type ProductImageResponse = {
  items?: ProductImageItem[];
};

type AssetReadUrlResponse = {
  url?: string;
};

type GoFrameResponse<T> = {
  code?: number;
  message?: string;
  data?: T;
};

type SellerApplicationShopDraft = {
  shop_no?: string;
  shopNo?: string;
  shop_name?: string;
  shopName?: string;
  shop_display_name?: string;
  shopDisplayName?: string;
};

type SellerApplicationItem = {
  status?: number | string;
  shop_submitted?: SellerApplicationShopDraft;
  shopSubmitted?: SellerApplicationShopDraft;
  shop_draft?: SellerApplicationShopDraft;
  shopDraft?: SellerApplicationShopDraft;
};

type ListSellerApplicationsResponse = {
  applications?: SellerApplicationItem[];
};

type ProductStoreCategoryBindingItem = {
  spu_no?: string;
  spuNo?: string;
  store_category_id?: number | string;
  storeCategoryId?: number | string;
  store_category_path?: Array<number | string>;
  storeCategoryPath?: Array<number | string>;
};

type BatchProductStoreCategoryBindingsResponse = {
  items?: ProductStoreCategoryBindingItem[];
};

type NormalizedSearchItem = {
  spuNo: string;
  title: string;
  coverAssetId: string;
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

type SearchShopPreviewItem = {
  spuNo: string;
  title: string;
  coverUrl: string;
  minPrice: number;
  salesCount: number;
};

type SearchShopCard = {
  shopNo: string;
  shopName: string;
  shopDisplayName: string;
  coverUrl: string;
  matchedProductCount: number;
  sampleSpuNo: string;
  previewItems: SearchShopPreviewItem[];
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

function unwrapGoFrame<T>(payload: unknown): T {
  if (!payload || typeof payload !== "object") {
    return payload as T;
  }

  const response = payload as GoFrameResponse<T>;
  if (typeof response.code === "number" && "data" in response) {
    if (response.code !== 0) {
      throw new Error(response.message || "upstream request failed");
    }
    return response.data as T;
  }
  return payload as T;
}

function normalizeSearchSvcItem(item: SearchSvcItem): NormalizedSearchItem {
  const categoryId = toNumber(item.categoryId ?? item.category_id, 0) || undefined;
  const descriptor = resolveBuyerCategoryDescriptor({
    categoryId,
    categoryName: toString(item.categoryName ?? item.category_name),
    title: toString(item.title)
  });

  return {
    spuNo: toString(item.spuNo ?? item.spu_no),
    title: toString(item.title),
    coverAssetId: toString(item.coverAssetId ?? item.cover_asset_id),
    coverUrl: toString(item.coverUrl ?? item.cover_url),
    minPrice: toNumber(item.minPrice ?? item.min_price),
    maxPrice: toNumber(item.maxPrice ?? item.max_price),
    categoryId,
    categoryName: descriptor.labelZh,
    shopNo: toString(item.shopNo ?? item.shop_no),
    shopName: toString(item.shopName ?? item.shop_name),
    salesCount: toNumber(item.salesCount ?? item.sales_count),
    stockTotal: toNumber(item.stockTotal ?? item.stock_total),
    avgScoreX100: toNumber(item.avgScoreX100 ?? item.avg_score_x100),
    reviewTotal: toNumber(item.reviewTotal ?? item.review_total)
  };
}

function normalizeCatalogItem(item: CatalogItem, imageUrl: string) {
  const categoryId = toNumber(item.categoryId ?? item.category_id, 0) || undefined;
  const descriptor = resolveBuyerCategoryDescriptor({
    categoryId,
    categoryName: toString(item.categoryName ?? item.category_name),
    title: toString(item.title)
  });

  return {
    spuNo: toString(item.spuNo ?? item.spu_no),
    title: toString(item.title),
    coverAssetId: primaryAssetIdOf(item),
    coverUrl: imageUrl,
    minPrice: toNumber(item.minSalePrice ?? item.min_sale_price),
    maxPrice: toNumber(item.maxSalePrice ?? item.max_sale_price),
    categoryId,
    categoryName: descriptor.labelZh,
    shopNo: "",
    shopName: "",
    salesCount: toNumber(item.soldCount ?? item.sold_count),
    stockTotal: 0,
    avgScoreX100: 0,
    reviewTotal: 0
  };
}

function toShopPreviewItem(item: Pick<NormalizedSearchItem, "spuNo" | "title" | "coverUrl" | "minPrice" | "salesCount">): SearchShopPreviewItem {
  return {
    spuNo: item.spuNo,
    title: item.title,
    coverUrl: item.coverUrl,
    minPrice: item.minPrice,
    salesCount: item.salesCount
  };
}

function primaryAssetIdOfSearchItem(item: NormalizedSearchItem): string {
  return item.coverAssetId.trim();
}

function primaryAssetIdOf(item: CatalogItem): string {
  const mainImageAssetIds = item.mainImageAssetIds ?? item.main_image_asset_ids;
  if (Array.isArray(mainImageAssetIds) && mainImageAssetIds.length > 0) {
    return toString(mainImageAssetIds[0]);
  }
  return toString(item.coverImageAssetId ?? item.cover_image_asset_id);
}

function matchesKeyword(item: CatalogItem, keyword: string) {
  const normalized = normalizeKeyword(keyword);
  if (!normalized) {
    return true;
  }

  const title = normalizeKeyword(toString(item.title));
  const subTitle = normalizeKeyword(toString(item.subTitle ?? item.sub_title));
  const spuNo = normalizeKeyword(toString(item.spuNo ?? item.spu_no));
  return title.includes(normalized) || subTitle.includes(normalized) || spuNo.includes(normalized);
}

function normalizeKeyword(value: string): string {
  return value.trim().toLowerCase();
}

function buildShopCards(query: string, items: NormalizedSearchItem[]): SearchShopCard[] {
  const keyword = normalizeKeyword(query);
  if (!keyword) {
    return [];
  }

  const shopMap = new Map<string, SearchShopCard>();
  for (const item of items) {
    const shopNo = item.shopNo.trim();
    const shopName = item.shopName.trim();
    if (!shopNo || !shopName) {
      continue;
    }

    const matched = normalizeKeyword(shopName).includes(keyword) || normalizeKeyword(shopNo).includes(keyword);
    if (!matched) {
      continue;
    }

    const existing = shopMap.get(shopNo);
    if (existing) {
      existing.matchedProductCount += 1;
      if (!existing.coverUrl && item.coverUrl) {
        existing.coverUrl = item.coverUrl;
      }
      if (!existing.sampleSpuNo && item.spuNo) {
        existing.sampleSpuNo = item.spuNo;
      }
      existing.previewItems = mergePreviewItems(existing.previewItems, [toShopPreviewItem(item)], 5);
      continue;
    }

    shopMap.set(shopNo, {
      shopNo,
      shopName,
      shopDisplayName: shopName,
      coverUrl: item.coverUrl,
      matchedProductCount: 1,
      sampleSpuNo: item.spuNo,
      previewItems: item.spuNo ? [toShopPreviewItem(item)] : []
    });
  }

  return Array.from(shopMap.values()).sort((left, right) => right.matchedProductCount - left.matchedProductCount);
}

async function fetchBuyerVisibleShopCards(query: string, limit: number): Promise<SearchShopCard[]> {
  const keyword = query.trim();
  if (!keyword) {
    return [];
  }

  const baseUrl = process.env.SELLER_SHOP_SERVICE_BASE_URL || "http://127.0.0.1:8007";
  const params = new URLSearchParams({
    keyword,
    page: "1",
    pageSize: String(limit)
  });
  const response = await fetch(`${baseUrl}/v1/admin/seller/applications?${params.toString()}`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    return [];
  }

  const payload = unwrapGoFrame<ListSellerApplicationsResponse>(await response.json());
  const shopMap = new Map<string, SearchShopCard>();
  for (const application of payload?.applications ?? []) {
    if (toNumber(application.status) !== 4) {
      continue;
    }

    const shop = application.shopSubmitted ?? application.shop_submitted ?? application.shopDraft ?? application.shop_draft;
    const shopNo = toString(shop?.shopNo ?? shop?.shop_no);
    const shopName = toString(shop?.shopName ?? shop?.shop_name);
    const shopDisplayName = toString(shop?.shopDisplayName ?? shop?.shop_display_name);
    if (!shopNo || !shopName) {
      continue;
    }
    if (!shopMap.has(shopNo)) {
      shopMap.set(shopNo, {
        shopNo,
        shopName,
        shopDisplayName: shopDisplayName || shopName,
        coverUrl: "",
        matchedProductCount: 0,
        sampleSpuNo: "",
        previewItems: []
      });
    }
  }
  return Array.from(shopMap.values());
}

function mergePreviewItems(primary: SearchShopPreviewItem[], secondary: SearchShopPreviewItem[], limit = 5): SearchShopPreviewItem[] {
  const merged = new Map<string, SearchShopPreviewItem>();
  for (const item of [...primary, ...secondary]) {
    if (!item.spuNo || merged.has(item.spuNo)) {
      continue;
    }
    merged.set(item.spuNo, item);
    if (merged.size >= limit) {
      break;
    }
  }
  return Array.from(merged.values());
}

function mergeShopCards(primary: SearchShopCard[], fallback: SearchShopCard[]): SearchShopCard[] {
  const merged = new Map<string, SearchShopCard>();
  for (const shop of [...primary, ...fallback]) {
    if (!shop.shopNo) {
      continue;
    }
    const existing = merged.get(shop.shopNo);
    if (!existing) {
      merged.set(shop.shopNo, { ...shop });
      continue;
    }
    if (!existing.shopDisplayName && shop.shopDisplayName) {
      existing.shopDisplayName = shop.shopDisplayName;
    }
    if (!existing.shopName && shop.shopName) {
      existing.shopName = shop.shopName;
    }
    if (!existing.coverUrl && shop.coverUrl) {
      existing.coverUrl = shop.coverUrl;
    }
    if (!existing.sampleSpuNo && shop.sampleSpuNo) {
      existing.sampleSpuNo = shop.sampleSpuNo;
    }
    existing.previewItems = mergePreviewItems(existing.previewItems, shop.previewItems, 5);
    if (shop.matchedProductCount > existing.matchedProductCount) {
      existing.matchedProductCount = shop.matchedProductCount;
    }
  }
  return Array.from(merged.values());
}

function buildSearchPreviewMap(items: NormalizedSearchItem[]) {
  const previewMap = new Map<string, SearchShopPreviewItem[]>();
  for (const item of items) {
    const shopNo = item.shopNo.trim();
    if (!shopNo || !item.spuNo) {
      continue;
    }
    const current = previewMap.get(shopNo) ?? [];
    previewMap.set(shopNo, mergePreviewItems(current, [toShopPreviewItem(item)], 5));
  }
  return previewMap;
}

async function listProductStoreCategoryBindings(shopNo: string, spuNos: string[]) {
  const normalizedSpuNos = Array.from(new Set(spuNos.map((item) => item.trim()).filter(Boolean)));
  const bindings = new Map<string, { storeCategoryId: number; storeCategoryPath: number[] }>();
  if (!shopNo || normalizedSpuNos.length === 0) {
    return bindings;
  }

  const baseUrl = process.env.SELLER_SHOP_SERVICE_BASE_URL || "http://127.0.0.1:8007";
  const response = await fetch(`${baseUrl}/v1/seller/shops/${encodeURIComponent(shopNo)}/products/store-categories:batch-get`, {
    method: "POST",
    cache: "no-store",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      spuNos: normalizedSpuNos
    })
  }).catch(() => null);
  if (!response?.ok) {
    return bindings;
  }

  const payload = unwrapGoFrame<BatchProductStoreCategoryBindingsResponse>(await response.json());
  for (const item of payload?.items ?? []) {
    const spuNo = toString(item.spuNo ?? item.spu_no);
    if (!spuNo) {
      continue;
    }
    const path = (item.storeCategoryPath ?? item.store_category_path ?? []).map((entry) => toNumber(entry)).filter((entry) => entry > 0);
    bindings.set(spuNo, {
      storeCategoryId: toNumber(item.storeCategoryId ?? item.store_category_id),
      storeCategoryPath: path
    });
  }

  return bindings;
}

async function filterCatalogItemsByStoreCategory(
  rawItems: CatalogItem[],
  shopNo: string,
  storeCategoryId: number,
  storeCategoryLevel: number
) {
  if (!shopNo || storeCategoryId <= 0 || rawItems.length === 0) {
    return rawItems;
  }

  const bindings = await listProductStoreCategoryBindings(
    shopNo,
    rawItems.map((item) => toString(item.spuNo ?? item.spu_no))
  );

  return rawItems.filter((item) => {
    const spuNo = toString(item.spuNo ?? item.spu_no);
    const binding = bindings.get(spuNo);
    if (!binding) {
      return false;
    }
    if (storeCategoryLevel === 2) {
      return binding.storeCategoryId === storeCategoryId;
    }
    return binding.storeCategoryPath.includes(storeCategoryId);
  });
}

async function trySearchSvc(
  query: string,
  pageSize: number,
  nextCursor: string,
  shopNo: string,
  storeCategoryId: number,
  storeCategoryLevel: number
) {
  const baseUrl = process.env.SEARCH_SERVICE_BASE_URL || "http://127.0.0.1:8018";
  const params = new URLSearchParams();
  params.set("query", query);
  params.set("pageSize", String(pageSize));
  if (nextCursor) {
    params.set("nextCursor", nextCursor);
  }
  if (shopNo) {
    params.set("shopNo", shopNo);
  }
  if (storeCategoryId > 0) {
    params.set("storeCategoryId", String(storeCategoryId));
    if (storeCategoryLevel > 0) {
      params.set("storeCategoryLevel", String(storeCategoryLevel));
    }
  }

  const response = await fetch(`${baseUrl}/v1/search/products?${params.toString()}`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    throw new Error(`search-svc request failed: ${response.status}`);
  }

  const payload = unwrapGoFrame<SearchSvcResponse>(await response.json());
  const rawItems = Array.isArray(payload?.list) ? payload.list.map((item) => normalizeSearchSvcItem(item)) : [];
  const assetUrlMap = await listMediaAssetReadUrls(rawItems.map((item) => primaryAssetIdOfSearchItem(item)));
  const items = rawItems.map((item) => ({
    ...item,
    coverUrl: assetUrlMap.get(primaryAssetIdOfSearchItem(item)) || item.coverUrl
  }));
  return {
    items,
    shops: buildShopCards(query, items),
    nextCursor: toString(payload?.nextCursor ?? payload?.next_cursor),
    hasMore: Boolean(payload?.hasMore ?? payload?.has_more),
    source: "search-svc" as const
  };
}

async function listCatalogProductImages(baseUrl: string, spuNos: string[]) {
  const normalized = spuNos.filter(Boolean);
  if (normalized.length === 0) {
    return new Map<string, string>();
  }

  const params = new URLSearchParams({
    spu_nos: normalized.join(",")
  });
  const response = await fetch(`${baseUrl}/v1/catalog/buyer/product-images?${params.toString()}`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    return new Map<string, string>();
  }

  const payload = unwrapGoFrame<ProductImageResponse>(await response.json());
  const imageMap = new Map<string, string>();
  for (const item of payload?.items ?? []) {
    const spuNo = toString(item.spuNo ?? item.spu_no);
    const imageUrl = toString(item.imageUrl ?? item.image_url);
    if (spuNo && imageUrl) {
      imageMap.set(spuNo, imageUrl);
    }
  }
  return imageMap;
}

async function listMediaAssetReadUrls(assetIds: string[]) {
  const normalized = Array.from(new Set(assetIds.map((item) => item.trim()).filter(Boolean)));
  if (normalized.length === 0) {
    return new Map<string, string>();
  }

  const baseUrl = process.env.MEDIA_SERVICE_BASE_URL || "http://127.0.0.1:8006";
  const entries = await Promise.all(
    normalized.map(async (assetId) => {
      const response = await fetch(
        `${baseUrl}/v1/media/assets/${encodeURIComponent(assetId)}/read-url?ttl_seconds=900`,
        {
          method: "GET",
          cache: "no-store"
        }
      ).catch(() => null);
      if (!response?.ok) {
        return [assetId, ""] as const;
      }

      const payload = unwrapGoFrame<AssetReadUrlResponse>(await response.json());
      return [assetId, toString(payload?.url)] as const;
    })
  );

  return new Map(entries.filter((entry) => entry[1]));
}

async function fetchCatalogItems(baseUrl: string, endpoint: string, pageSize: number, query: string) {
  const params = new URLSearchParams({
    page: "1",
    page_size: String(pageSize),
    sort_by: "4"
  });
  if (query) {
    params.set("keyword", query);
  }

  const response = await fetch(`${baseUrl}${endpoint}?${params.toString()}`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    return [];
  }

  const payload = unwrapGoFrame<CatalogSearchResponse>(await response.json());
  const list = payload?.items ?? payload?.products;
  return Array.isArray(list) ? list : [];
}

async function fetchSellerCatalogItemsPage(params: {
  shopNo: string;
  keyword: string;
  page: number;
  pageSize: number;
  statuses?: number[];
}) {
  const baseUrl = process.env.CATALOG_SERVICE_BASE_URL || "http://127.0.0.1:8004";
  const query = new URLSearchParams({
    page: String(params.page),
    page_size: String(params.pageSize)
  });
  if (params.keyword.trim()) {
    query.set("keyword", params.keyword.trim());
  }
  for (const status of params.statuses ?? []) {
    query.append("statuses", String(status));
  }

  const response = await fetch(`${baseUrl}/v1/catalog/seller/products?${query.toString()}`, {
    method: "GET",
    cache: "no-store"
  }).catch(() => null);
  if (!response?.ok) {
    return { items: [] as CatalogItem[], total: 0 };
  }

  const payload = unwrapGoFrame<CatalogSearchResponse>(await response.json());
  const list = payload?.products ?? payload?.items ?? [];
  const items = Array.isArray(list)
    ? list.filter((item) => {
        const currentShopNo = toString(item.shopNo ?? item.shop_no);
        const currentStatus = toNumber(item.spuStatus ?? item.spu_status);
        return currentShopNo === params.shopNo && currentStatus === 4 && matchesKeyword(item, params.keyword);
      })
    : [];

  return {
    items,
    total: toNumber(payload?.total, 0)
  };
}

async function listSellerCatalogItemsForShop(shopNo: string, keyword: string, desiredCount: number) {
  if (!shopNo) {
    return [] as CatalogItem[];
  }

  const pageSize = Math.min(Math.max(desiredCount * 2, 80), 200);
  const maxPages = 10;
  const matched = new Map<string, CatalogItem>();
  let total = Number.POSITIVE_INFINITY;

  for (let page = 1; page <= maxPages; page += 1) {
    const result = await fetchSellerCatalogItemsPage({
      shopNo,
      keyword,
      page,
      pageSize,
      statuses: [4]
    });

    if (!Number.isFinite(total)) {
      total = result.total;
    }

    for (const item of result.items) {
      const spuNo = toString(item.spuNo ?? item.spu_no);
      if (!spuNo || matched.has(spuNo)) {
        continue;
      }
      matched.set(spuNo, item);
      if (matched.size >= desiredCount) {
        return Array.from(matched.values());
      }
    }

    if (result.total <= page * pageSize) {
      break;
    }
  }

  return Array.from(matched.values());
}

async function fetchShopPreviewItems(shopNo: string, limit: number) {
  if (!shopNo) {
    return [] as SearchShopPreviewItem[];
  }

  const previewSeed = (await listSellerCatalogItemsForShop(shopNo, "", Math.max(limit * 2, 12)))
    .sort((left, right) => toNumber(right.soldCount ?? right.sold_count) - toNumber(left.soldCount ?? left.sold_count))
    .slice(0, limit);

  const buyerApiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || "http://127.0.0.1:8000";
  const imageMap = await listCatalogProductImages(
    buyerApiBaseUrl,
    previewSeed.map((item) => toString(item.spuNo ?? item.spu_no))
  );
  const assetUrlMap = await listMediaAssetReadUrls(previewSeed.map((item) => primaryAssetIdOf(item)));

  return previewSeed.map((item) => {
    const spuNo = toString(item.spuNo ?? item.spu_no);
    return {
      spuNo,
      title: toString(item.title),
      coverUrl: imageMap.get(spuNo) ?? assetUrlMap.get(primaryAssetIdOf(item)) ?? "",
      minPrice: toNumber(item.minSalePrice ?? item.min_sale_price),
      salesCount: toNumber(item.soldCount ?? item.sold_count)
    };
  });
}

async function enrichShopCards(cards: SearchShopCard[], items: NormalizedSearchItem[]) {
  if (cards.length === 0) {
    return cards;
  }

  const searchPreviewMap = buildSearchPreviewMap(items);
  const enrichedCards = await Promise.all(
    cards.map(async (shop) => {
      const searchPreviews = searchPreviewMap.get(shop.shopNo) ?? [];
      const previewItems = mergePreviewItems(shop.previewItems, searchPreviews, 5);
      const catalogPreviews = previewItems.length >= 4 ? [] : await fetchShopPreviewItems(shop.shopNo, 5);
      const mergedPreviewItems = mergePreviewItems(previewItems, catalogPreviews, 5);
      const leadPreview = mergedPreviewItems[0];

      return {
        ...shop,
        coverUrl: shop.coverUrl || leadPreview?.coverUrl || "",
        sampleSpuNo: shop.sampleSpuNo || leadPreview?.spuNo || "",
        previewItems: mergedPreviewItems
      };
    })
  );

  return enrichedCards;
}

async function fallbackCatalog(query: string, pageSize: number) {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || "http://127.0.0.1:8000";
  const keyword = query.trim().toLowerCase();
  let rawItems = await fetchCatalogItems(
    baseUrl,
    query ? "/v1/catalog/buyer/products/search" : "/v1/catalog/buyer/products",
    pageSize,
    query
  );

  if (query && rawItems.length === 0) {
    const fallbackPool = await fetchCatalogItems(baseUrl, "/v1/catalog/buyer/products", Math.max(pageSize, 100), "");
    rawItems = fallbackPool.filter((item) => {
      const spuNo = toString(item.spuNo ?? item.spu_no).toLowerCase();
      const title = toString(item.title).toLowerCase();
      return spuNo.includes(keyword) || title.includes(keyword);
    });
  }

  const imageMap = await listCatalogProductImages(
    baseUrl,
    rawItems.map((item) => toString(item.spuNo ?? item.spu_no))
  );
  const assetUrlMap = await listMediaAssetReadUrls(rawItems.map((item) => primaryAssetIdOf(item)));
  const items = rawItems.map((item) => {
    const spuNo = toString(item.spuNo ?? item.spu_no);
    return normalizeCatalogItem(item, imageMap.get(spuNo) ?? assetUrlMap.get(primaryAssetIdOf(item)) ?? "");
  });
  const shops = await fetchBuyerVisibleShopCards(query, 10);

  return {
    items,
    shops: await enrichShopCards(shops, items),
    nextCursor: "",
    hasMore: false,
    source: "catalog-fallback" as const
  };
}

async function fallbackShopCatalog(
  query: string,
  pageSize: number,
  shopNo: string,
  storeCategoryId: number,
  storeCategoryLevel: number
) {
  const rawItems = await listSellerCatalogItemsForShop(shopNo, query, Math.max(pageSize * 4, 120));
  const filteredItems = await filterCatalogItemsByStoreCategory(rawItems, shopNo, storeCategoryId, storeCategoryLevel);

  const imageMap = await listCatalogProductImages(
    process.env.NEXT_PUBLIC_API_BASE_URL || "http://127.0.0.1:8000",
    filteredItems.map((item) => toString(item.spuNo ?? item.spu_no))
  );
  const assetUrlMap = await listMediaAssetReadUrls(filteredItems.map((item) => primaryAssetIdOf(item)));

  return {
    items: filteredItems.map((item) => {
      const spuNo = toString(item.spuNo ?? item.spu_no);
      return {
        ...normalizeCatalogItem(item, imageMap.get(spuNo) ?? assetUrlMap.get(primaryAssetIdOf(item)) ?? ""),
        shopNo
      };
    }),
    shops: [],
    nextCursor: "",
    hasMore: false,
    source: "catalog-fallback" as const
  };
}

export async function GET(request: NextRequest) {
  const query = request.nextUrl.searchParams.get("q")?.trim() || "";
  const pageSize = Math.min(Math.max(Number(request.nextUrl.searchParams.get("pageSize") || "24"), 1), 100);
  const nextCursor = request.nextUrl.searchParams.get("nextCursor")?.trim() || "";
  const shopNo = request.nextUrl.searchParams.get("shopNo")?.trim() || "";
  const storeCategoryId = Math.max(Number(request.nextUrl.searchParams.get("storeCategoryId") || "0"), 0);
  const storeCategoryLevel = Math.max(Number(request.nextUrl.searchParams.get("storeCategoryLevel") || "0"), 0);

  try {
    const result = await trySearchSvc(query, pageSize, nextCursor, shopNo, storeCategoryId, storeCategoryLevel);
    if (shopNo) {
      if (result.items.length > 0) {
        return NextResponse.json(result);
      }
      return NextResponse.json(await fallbackShopCatalog(query, pageSize, shopNo, storeCategoryId, storeCategoryLevel));
    }
    if (!shopNo && query) {
      result.shops = mergeShopCards(result.shops, await fetchBuyerVisibleShopCards(query, 10));
      result.shops = await enrichShopCards(result.shops, result.items);
    }
    if (result.items.length > 0 || !query) {
      return NextResponse.json(result);
    }
  } catch {
    if (shopNo) {
      return NextResponse.json(await fallbackShopCatalog(query, pageSize, shopNo, storeCategoryId, storeCategoryLevel));
    }
  }

  if (shopNo) {
    return NextResponse.json(await fallbackShopCatalog(query, pageSize, shopNo, storeCategoryId, storeCategoryLevel));
  }

  const fallback = await fallbackCatalog(query, pageSize);
  return NextResponse.json(fallback);
}
