import { headers } from "next/headers";
import ShopPageClient from "./shop-page-client";
import type { SearchProductsResult } from "@/features/search/types";
import type { MallShopDetail, ShopStoreCategory } from "@/features/shop/types";

function inferIsZh() {
  const headerStore = headers();
  const acceptLanguage = headerStore.get("accept-language") || "";
  return acceptLanguage.toLowerCase().includes("zh");
}

function resolveOrigin() {
  const headerStore = headers();
  const host = headerStore.get("x-forwarded-host") || headerStore.get("host") || "127.0.0.1:3100";
  const protocol = headerStore.get("x-forwarded-proto") || (host.includes("localhost") || host.includes("127.0.0.1") ? "http" : "https");
  return `${protocol}://${host}`;
}

async function fetchJson<T>(url: string): Promise<T> {
  const response = await fetch(url, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    throw new Error(`request failed: ${response.status}`);
  }
  return (await response.json()) as T;
}

async function loadInitialShop(origin: string, shopNo: string): Promise<MallShopDetail | null> {
  try {
    return await fetchJson<MallShopDetail>(`${origin}/api/shops/${encodeURIComponent(shopNo)}`);
  } catch {
    return null;
  }
}

async function loadInitialStoreCategories(origin: string, shopNo: string): Promise<ShopStoreCategory[]> {
  try {
    return await fetchJson<ShopStoreCategory[]>(`${origin}/api/shops/${encodeURIComponent(shopNo)}/store-categories`);
  } catch {
    return [];
  }
}

async function loadInitialProducts(
  origin: string,
  shopNo: string,
  keyword: string,
  storeCategoryId?: number,
  storeCategoryLevel?: 1 | 2
): Promise<SearchProductsResult> {
  const params = new URLSearchParams();
  params.set("shopNo", shopNo);
  params.set("pageSize", "80");
  if (keyword.trim()) {
    params.set("q", keyword.trim());
  }
  if (storeCategoryId && storeCategoryId > 0) {
    params.set("storeCategoryId", String(storeCategoryId));
    if (storeCategoryLevel) {
      params.set("storeCategoryLevel", String(storeCategoryLevel));
    }
  }

  try {
    return await fetchJson<SearchProductsResult>(`${origin}/api/search/products?${params.toString()}`);
  } catch {
    return {
      items: [],
      shops: [],
      nextCursor: "",
      hasMore: false,
      source: "catalog-fallback"
    };
  }
}

export default async function ShopPage({
  params,
  searchParams
}: {
  params: { shopNo: string };
  searchParams?: { q?: string; from?: string; primary?: string; secondary?: string; sort?: string };
}) {
  const allowedSorts = new Set(["default", "sales", "price-asc", "price-desc"]);
  const shopNo = params.shopNo?.trim() || "";
  const keyword = searchParams?.q?.trim() || "";
  const fromSearch = searchParams?.from === "search";
  const primaryCategoryId = Math.max(Number(searchParams?.primary?.trim() || "0"), 0);
  const secondaryCategoryId = Math.max(Number(searchParams?.secondary?.trim() || "0"), 0);
  const requestedSort = searchParams?.sort?.trim() || "default";
  const selectedSort = (allowedSorts.has(requestedSort) ? requestedSort : "default") as "default" | "sales" | "price-asc" | "price-desc";
  const origin = resolveOrigin();
  const isZh = inferIsZh();
  const selectedStoreCategoryId = secondaryCategoryId || primaryCategoryId;
  const selectedStoreCategoryLevel = secondaryCategoryId ? 2 : primaryCategoryId ? 1 : undefined;

  const [initialShop, initialCategories, initialProducts] = await Promise.all([
    loadInitialShop(origin, shopNo),
    loadInitialStoreCategories(origin, shopNo),
    loadInitialProducts(origin, shopNo, keyword, selectedStoreCategoryId, selectedStoreCategoryLevel as 1 | 2 | undefined)
  ]);

  return (
    <ShopPageClient
      shopNo={shopNo}
      keyword={keyword}
      fromSearch={fromSearch}
      isZh={isZh}
      selectedPrimaryCategoryId={primaryCategoryId}
      selectedSecondaryCategoryId={secondaryCategoryId}
      selectedSort={selectedSort}
      initialShop={initialShop}
      initialCategories={initialCategories}
      initialProducts={initialProducts}
    />
  );
}
