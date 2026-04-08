import { apiClient } from "@/lib/api-client";
import {
  BUYER_SORT_BY,
  BuyerProductDetail,
  BuyerProductDetailSku,
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

type RawProductSku = {
  sku_no?: string;
  skuNo?: string;
  spu_no?: string;
  spuNo?: string;
  shop_no?: string;
  shopNo?: string;
  sku_name?: string;
  skuName?: string;
  sku_image_asset_id?: number | string;
  skuImageAssetId?: number | string;
  sale_price?: number | string;
  salePrice?: number | string;
  market_price?: number | string;
  marketPrice?: number | string;
  stock_status?: number | string;
  stockStatus?: number | string;
  sale_attrs_json?: string;
  saleAttrsJson?: string;
};

type RawProductSpu = {
  spu_no?: string;
  spuNo?: string;
  shop_no?: string;
  shopNo?: string;
  title?: string;
  sub_title?: string;
  subTitle?: string;
  min_sale_price?: number | string;
  minSalePrice?: number | string;
  max_sale_price?: number | string;
  maxSalePrice?: number | string;
  min_market_price?: number | string;
  minMarketPrice?: number | string;
  max_market_price?: number | string;
  maxMarketPrice?: number | string;
  main_image_asset_ids?: Array<number | string>;
  mainImageAssetIds?: Array<number | string>;
  detail_image_asset_ids?: Array<number | string>;
  detailImageAssetIds?: Array<number | string>;
};

type RawProductAggregate = {
  spu?: RawProductSpu;
  skus?: RawProductSku[];
};

type RawGetProductDetailRes = {
  product?: RawProductAggregate;
};

function toNumber(value: unknown, fallback = 0): number {
  const num = Number(value);
  return Number.isFinite(num) ? num : fallback;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "number" || typeof value === "boolean" || typeof value === "bigint") {
    return String(value);
  }
  return "";
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

function normalizeSku(raw: RawProductSku): BuyerProductDetailSku {
  const saleAttrs = raw.sale_attrs_json ?? raw.saleAttrsJson;
  return {
    skuNo: toString(raw.sku_no ?? raw.skuNo),
    spuNo: toString(raw.spu_no ?? raw.spuNo),
    shopNo: toString(raw.shop_no ?? raw.shopNo),
    skuName: toString(raw.sku_name ?? raw.skuName),
    skuImageAssetId: toString(raw.sku_image_asset_id ?? raw.skuImageAssetId),
    salePrice: toNumber(raw.sale_price ?? raw.salePrice, 0),
    marketPrice: toNumber(raw.market_price ?? raw.marketPrice, 0),
    stockStatus: toNumber(raw.stock_status ?? raw.stockStatus, 0),
    saleAttrsJson: toString((saleAttrs ?? "[]"))
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

export async function listBuyerAssetReadUrls(assetIds: string[], ttlSeconds = 900): Promise<Record<string, string>> {
  const normalized = Array.from(new Set(assetIds.map((item) => item.trim()).filter((item) => item.length > 0))).slice(0, 50);
  if (normalized.length === 0) {
    return {};
  }

  const results = await Promise.all(
    normalized.map(async (assetId) => {
      try {
        const response = await apiClient.get<{ url?: string }>(`/v1/media/assets/${encodeURIComponent(assetId)}/read-url`, {
          params: {
            ttl_seconds: ttlSeconds
          }
        });
        return {
          assetId,
          url: toString(response?.url)
        };
      } catch {
        return {
          assetId,
          url: ""
        };
      }
    })
  );

  return results.reduce<Record<string, string>>((acc, item) => {
    if (!item.assetId || !item.url || acc[item.assetId]) {
      return acc;
    }
    acc[item.assetId] = item.url;
    return acc;
  }, {});
}

export async function getBuyerProductDetail(spuNo: string): Promise<BuyerProductDetail> {
  const normalizedSpuNo = spuNo.trim();
  if (!normalizedSpuNo) {
    throw new Error("spu_no is required");
  }

  const response = await apiClient.get<RawGetProductDetailRes>(`/v1/catalog/buyer/products/${encodeURIComponent(normalizedSpuNo)}`);
  const spu = response?.product?.spu ?? {};
  const skus = Array.isArray(response?.product?.skus) ? response.product!.skus!.map((item) => normalizeSku(item)) : [];

  return {
    spuNo: toString(spu.spu_no ?? spu.spuNo ?? normalizedSpuNo),
    shopNo: toString(spu.shop_no ?? spu.shopNo),
    title: toString(spu.title),
    subTitle: toString(spu.sub_title ?? spu.subTitle),
    minSalePrice: toNumber(spu.min_sale_price ?? spu.minSalePrice, 0),
    maxSalePrice: toNumber(spu.max_sale_price ?? spu.maxSalePrice, 0),
    minMarketPrice: toNumber(spu.min_market_price ?? spu.minMarketPrice, 0),
    maxMarketPrice: toNumber(spu.max_market_price ?? spu.maxMarketPrice, 0),
    mainImageAssetIds: (spu.main_image_asset_ids ?? spu.mainImageAssetIds ?? []).map((item) => String(item)),
    detailImageAssetIds: (spu.detail_image_asset_ids ?? spu.detailImageAssetIds ?? []).map((item) => String(item)),
    skus
  };
}
