import { apiClient } from "@/lib/api-client";
import {
  BUYER_SORT_BY,
  BuyerProductCard,
  BuyerProductImageItem,
  ListBuyerProductsParams,
  ListBuyerProductsResult
} from "./types";

type RawBuyerProductCard = {
  spu_no?: string;
  spuNo?: string;
  title?: string;
  category_id?: number | string;
  categoryId?: number | string;
  brand_no?: string;
  brandNo?: string;
  cover_image_asset_id?: number | string;
  coverImageAssetId?: number | string;
  min_sale_price?: number | string;
  minSalePrice?: number | string;
  max_sale_price?: number | string;
  maxSalePrice?: number | string;
  sold_count?: number | string;
  soldCount?: number | string;
};

type RawListBuyerProductsResponse = {
  items?: RawBuyerProductCard[];
  page?: number | string;
  page_size?: number | string;
  pageSize?: number | string;
  total?: number | string;
};

type RawBuyerProductImageItem = {
  spu_no?: string;
  spuNo?: string;
  image_url?: string;
  imageUrl?: string;
};

type RawListBuyerProductImagesResponse = {
  items?: RawBuyerProductImageItem[];
};

function toNumber(value: unknown, fallback = 0): number {
  const num = Number(value);
  return Number.isFinite(num) ? num : fallback;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function normalizeCard(raw: RawBuyerProductCard): BuyerProductCard {
  return {
    spuNo: toString(raw.spu_no ?? raw.spuNo),
    title: toString(raw.title),
    categoryId: toNumber(raw.category_id ?? raw.categoryId, 0),
    brandNo: toString(raw.brand_no ?? raw.brandNo),
    coverImageAssetId: toString(raw.cover_image_asset_id ?? raw.coverImageAssetId),
    minSalePrice: toNumber(raw.min_sale_price ?? raw.minSalePrice, 0),
    maxSalePrice: toNumber(raw.max_sale_price ?? raw.maxSalePrice, 0),
    soldCount: toNumber(raw.sold_count ?? raw.soldCount, 0)
  };
}

export async function listBuyerProducts(params: ListBuyerProductsParams = {}): Promise<ListBuyerProductsResult> {
  const requestParams: Record<string, number> = {
    page: params.page ?? 1,
    page_size: params.pageSize ?? 24,
    sort_by: params.sortBy ?? BUYER_SORT_BY.SALES_DESC
  };
  if (typeof params.categoryId === "number") {
    requestParams.category_id = params.categoryId;
  }

  const response = await apiClient.get<RawListBuyerProductsResponse>("/v1/catalog/buyer/products", {
    params: requestParams
  });

  const items = Array.isArray(response?.items) ? response.items.map((item) => normalizeCard(item)) : [];

  return {
    items,
    page: toNumber(response?.page, 1),
    pageSize: toNumber(response?.page_size ?? response?.pageSize, params.pageSize ?? 24),
    total: toNumber(response?.total, items.length)
  };
}

export async function listBuyerProductImages(spuNos: string[]): Promise<Record<string, string>> {
  const normalized = Array.from(new Set(spuNos.map((item) => item.trim()).filter((item) => item.length > 0))).slice(0, 200);
  if (normalized.length === 0) {
    return {};
  }

  const response = await apiClient.get<RawListBuyerProductImagesResponse>("/v1/catalog/buyer/product-images", {
    params: {
      spu_nos: normalized.join(",")
    }
  });

  const items: BuyerProductImageItem[] = Array.isArray(response?.items)
    ? response.items.map((row) => ({
        spuNo: toString(row.spu_no ?? row.spuNo),
        imageUrl: toString(row.image_url ?? row.imageUrl)
      }))
    : [];

  return items.reduce<Record<string, string>>((acc, item) => {
    if (!item.spuNo || !item.imageUrl || acc[item.spuNo]) {
      return acc;
    }
    acc[item.spuNo] = item.imageUrl;
    return acc;
  }, {});
}
