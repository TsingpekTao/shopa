import { apiClient } from "@/lib/api-client";
import {
  CreateProductDraftInput,
  CreateProductDraftResult,
  ListSellerProductsParams,
  ListSellerProductsResult,
  SellerProductDetail,
  SellerProductReviewInfo,
  SellerProductSku,
  SellerProductSpu,
  SpuStatusCode,
  StockStatusCode
} from "./types";

type TimestampLike = {
  seconds?: number | string;
  nanos?: number | string;
};

type RawReviewInfo = {
  reject_reason_code?: string;
  rejectReasonCode?: string;
  reject_comment?: string;
  rejectComment?: string;
  reviewer_id?: number | string;
  reviewerId?: number | string;
  reviewed_at?: string | TimestampLike;
  reviewedAt?: string | TimestampLike;
};

type RawProductSku = {
  sku_no?: string;
  skuNo?: string;
  spu_no?: string;
  spuNo?: string;
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
  stock_version?: number | string;
  stockVersion?: number | string;
  sort_order?: number;
  sortOrder?: number;
  created_at?: string | TimestampLike;
  createdAt?: string | TimestampLike;
  updated_at?: string | TimestampLike;
  updatedAt?: string | TimestampLike;
};

type RawProductSpu = {
  spu_no?: string;
  spuNo?: string;
  shop_no?: string;
  shopNo?: string;
  title?: string;
  sub_title?: string;
  subTitle?: string;
  category_id?: number | string;
  categoryId?: number | string;
  brand_no?: string;
  brandNo?: string;
  main_image_asset_ids?: Array<number | string>;
  mainImageAssetIds?: Array<number | string>;
  detail_image_asset_ids?: Array<number | string>;
  detailImageAssetIds?: Array<number | string>;
  spu_status?: number | string;
  spuStatus?: number | string;
  spu_stock_status?: number | string;
  spuStockStatus?: number | string;
  min_sale_price?: number | string;
  minSalePrice?: number | string;
  max_sale_price?: number | string;
  maxSalePrice?: number | string;
  min_market_price?: number | string;
  minMarketPrice?: number | string;
  max_market_price?: number | string;
  maxMarketPrice?: number | string;
  version?: number;
  created_at?: string | TimestampLike;
  createdAt?: string | TimestampLike;
  updated_at?: string | TimestampLike;
  updatedAt?: string | TimestampLike;
};

type RawProductAggregate = {
  spu?: RawProductSpu;
  skus?: RawProductSku[];
};

type RawListProductsResponse = {
  products?: RawProductSpu[];
  page?: number;
  page_size?: number;
  pageSize?: number;
  total?: number;
};

type RawGetProductResponse = {
  product?: RawProductAggregate;
  review?: RawReviewInfo;
};

type RawCreateDraftResponse = {
  product?: RawProductAggregate;
};

const spuStatusByNumber: Record<number, SpuStatusCode> = {
  0: "SPU_STATUS_UNSPECIFIED",
  1: "SPU_STATUS_DRAFT",
  2: "SPU_STATUS_REVIEWING",
  3: "SPU_STATUS_APPROVED",
  4: "SPU_STATUS_ON_SHELF",
  5: "SPU_STATUS_OFF_SHELF",
  6: "SPU_STATUS_REJECTED",
  7: "SPU_STATUS_FROZEN",
  8: "SPU_STATUS_DELETED"
};

const stockStatusByNumber: Record<number, StockStatusCode> = {
  0: "STOCK_STATUS_UNSPECIFIED",
  1: "STOCK_STATUS_IN_STOCK",
  2: "STOCK_STATUS_OUT_OF_STOCK"
};

function toNumber(value: unknown, fallback = 0): number {
  const n = Number(value);
  return Number.isFinite(n) ? n : fallback;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function toTimestamp(value: unknown): string | undefined {
  if (!value) {
    return undefined;
  }
  if (typeof value === "string") {
    return value;
  }
  if (typeof value !== "object") {
    return undefined;
  }
  const raw = value as TimestampLike;
  const seconds = toNumber(raw.seconds, 0);
  if (!seconds) {
    return undefined;
  }
  const nanos = toNumber(raw.nanos, 0);
  return new Date(seconds * 1000 + Math.floor(nanos / 1_000_000)).toISOString();
}

function toAssetIds(input: Array<number | string> | undefined): string[] {
  if (!Array.isArray(input)) {
    return [];
  }
  return input.map((item) => String(item));
}

function toSpuStatus(value: unknown): SpuStatusCode {
  if (typeof value === "number") {
    return spuStatusByNumber[value] ?? "SPU_STATUS_UNSPECIFIED";
  }
  if (typeof value === "string") {
    const trimmed = value.trim().toUpperCase();
    if (trimmed in (spuStatusByNumber as unknown as Record<string, unknown>)) {
      return spuStatusByNumber[toNumber(trimmed, 0)] ?? "SPU_STATUS_UNSPECIFIED";
    }
    if (trimmed.startsWith("SPU_STATUS_")) {
      return trimmed as SpuStatusCode;
    }
    const normalized = `SPU_STATUS_${trimmed}` as SpuStatusCode;
    return normalized;
  }
  return "SPU_STATUS_UNSPECIFIED";
}

function toStockStatus(value: unknown): StockStatusCode {
  if (typeof value === "number") {
    return stockStatusByNumber[value] ?? "STOCK_STATUS_UNSPECIFIED";
  }
  if (typeof value === "string") {
    const trimmed = value.trim().toUpperCase();
    if (trimmed in (stockStatusByNumber as unknown as Record<string, unknown>)) {
      return stockStatusByNumber[toNumber(trimmed, 0)] ?? "STOCK_STATUS_UNSPECIFIED";
    }
    if (trimmed.startsWith("STOCK_STATUS_")) {
      return trimmed as StockStatusCode;
    }
    const normalized = `STOCK_STATUS_${trimmed}` as StockStatusCode;
    return normalized;
  }
  return "STOCK_STATUS_UNSPECIFIED";
}

function normalizeSpu(raw: RawProductSpu | undefined): SellerProductSpu {
  return {
    spuNo: toString(raw?.spu_no ?? raw?.spuNo),
    shopNo: toString(raw?.shop_no ?? raw?.shopNo),
    title: toString(raw?.title),
    subTitle: toString(raw?.sub_title ?? raw?.subTitle),
    categoryId: toNumber(raw?.category_id ?? raw?.categoryId, 0),
    brandNo: toString(raw?.brand_no ?? raw?.brandNo),
    mainImageAssetIds: toAssetIds(raw?.main_image_asset_ids ?? raw?.mainImageAssetIds),
    detailImageAssetIds: toAssetIds(raw?.detail_image_asset_ids ?? raw?.detailImageAssetIds),
    spuStatus: toSpuStatus(raw?.spu_status ?? raw?.spuStatus),
    spuStockStatus: toStockStatus(raw?.spu_stock_status ?? raw?.spuStockStatus),
    minSalePrice: toNumber(raw?.min_sale_price ?? raw?.minSalePrice, 0),
    maxSalePrice: toNumber(raw?.max_sale_price ?? raw?.maxSalePrice, 0),
    minMarketPrice: toNumber(raw?.min_market_price ?? raw?.minMarketPrice, 0),
    maxMarketPrice: toNumber(raw?.max_market_price ?? raw?.maxMarketPrice, 0),
    version: toNumber(raw?.version, 0),
    createdAt: toTimestamp(raw?.created_at ?? raw?.createdAt),
    updatedAt: toTimestamp(raw?.updated_at ?? raw?.updatedAt)
  };
}

function normalizeSku(raw: RawProductSku): SellerProductSku {
  return {
    skuNo: toString(raw.sku_no ?? raw.skuNo),
    spuNo: toString(raw.spu_no ?? raw.spuNo),
    skuName: toString(raw.sku_name ?? raw.skuName),
    skuImageAssetId: toString(raw.sku_image_asset_id ?? raw.skuImageAssetId) || undefined,
    salePrice: toNumber(raw.sale_price ?? raw.salePrice, 0),
    marketPrice: toNumber(raw.market_price ?? raw.marketPrice, 0),
    stockStatus: toStockStatus(raw.stock_status ?? raw.stockStatus),
    stockVersion: toNumber(raw.stock_version ?? raw.stockVersion, 0),
    sortOrder: toNumber(raw.sort_order ?? raw.sortOrder, 0),
    createdAt: toTimestamp(raw.created_at ?? raw.createdAt),
    updatedAt: toTimestamp(raw.updated_at ?? raw.updatedAt)
  };
}

function normalizeReview(raw: RawReviewInfo | undefined): SellerProductReviewInfo | undefined {
  if (!raw) {
    return undefined;
  }
  const reviewedAt = toTimestamp(raw.reviewed_at ?? raw.reviewedAt);
  const rejectReasonCode = toString(raw.reject_reason_code ?? raw.rejectReasonCode);
  const rejectComment = toString(raw.reject_comment ?? raw.rejectComment);
  const reviewerId = toString(raw.reviewer_id ?? raw.reviewerId);
  if (!rejectReasonCode && !rejectComment && !reviewerId && !reviewedAt) {
    return undefined;
  }
  return {
    rejectReasonCode: rejectReasonCode || undefined,
    rejectComment: rejectComment || undefined,
    reviewerId: reviewerId || undefined,
    reviewedAt
  };
}

function buildListQuery(params?: ListSellerProductsParams): string {
  const page = params?.page ?? 1;
  const pageSize = params?.pageSize ?? 20;
  const query = new URLSearchParams();
  query.set("page", String(page));
  query.set("page_size", String(pageSize));
  if (params?.keyword?.trim()) {
    query.set("keyword", params.keyword.trim());
  }
  (params?.statuses ?? []).forEach((status) => {
    const mapped = toNumber(spuStatusToNumber(status), 0);
    if (mapped > 0) {
      query.append("statuses", String(mapped));
    }
  });
  return query.toString();
}

function spuStatusToNumber(status: SpuStatusCode): number {
  for (const [num, code] of Object.entries(spuStatusByNumber)) {
    if (code === status) {
      return Number(num);
    }
  }
  return 0;
}

export async function listSellerProducts(params?: ListSellerProductsParams): Promise<ListSellerProductsResult> {
  const page = params?.page ?? 1;
  const pageSize = params?.pageSize ?? 20;
  const query = buildListQuery(params);
  const endpoint = query ? `/v1/catalog/seller/products?${query}` : "/v1/catalog/seller/products";
  const response = await apiClient.get<RawListProductsResponse>(endpoint);

  return {
    products: (response.products ?? []).map((item) => normalizeSpu(item)).filter((item) => item.spuNo),
    page: toNumber(response.page, page),
    pageSize: toNumber(response.page_size ?? response.pageSize, pageSize),
    total: toNumber(response.total, 0)
  };
}

export async function getSellerProductDetail(spuNo: string): Promise<SellerProductDetail> {
  const response = await apiClient.get<RawGetProductResponse>(`/v1/catalog/seller/products/${spuNo}`);
  const aggregate = response.product;
  return {
    spu: aggregate?.spu ? normalizeSpu(aggregate.spu) : null,
    skus: (aggregate?.skus ?? []).map((item) => normalizeSku(item)).filter((item) => item.skuNo),
    review: normalizeReview(response.review)
  };
}

export async function createProductDraft(payload: CreateProductDraftInput): Promise<CreateProductDraftResult> {
  const response = await apiClient.post<RawCreateDraftResponse>("/v1/catalog/seller/products/draft", {
    shop_no: payload.shopNo,
    title: payload.title,
    sub_title: payload.subTitle ?? "",
    category_id: payload.categoryId ?? 0,
    brand_no: payload.brandNo ?? "",
    main_image_asset_ids: (payload.mainImageAssetIds ?? []).map((item) => Number(item)).filter((item) => Number.isFinite(item)),
    detail_image_asset_ids: (payload.detailImageAssetIds ?? [])
      .map((item) => Number(item))
      .filter((item) => Number.isFinite(item)),
    spu_attrs: []
  });

  const spu = normalizeSpu(response.product?.spu);
  return {
    spuNo: spu.spuNo,
    version: spu.version
  };
}
