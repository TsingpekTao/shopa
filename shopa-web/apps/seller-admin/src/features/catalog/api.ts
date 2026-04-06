import { apiClient } from "@/lib/api-client";
import {
  CreateProductDraftInput,
  CreateProductDraftResult,
  ProductDetailResult,
  SaveDraftResult,
  ListSellerProductsParams,
  ListSellerProductsResult,
  SaveSellerProductDraftInput,
  DeleteProductDraftInput,
  spuStatusCodeMap,
  SellerCatalogDraftPayload,
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
  sale_attrs?: Array<{ attr_code?: string; attrCode?: string; value?: string }>;
  saleAttrs?: Array<{ attr_code?: string; attrCode?: string; value?: string }>;
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
  spu_attrs?: Array<{ attr_code?: string; attrCode?: string; value?: string }>;
  spuAttrs?: Array<{ attr_code?: string; attrCode?: string; value?: string }>;
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

function toMinorUnits(amount: number): number {
  return Math.max(0, Math.round((Number.isFinite(amount) ? amount : 0) * 100));
}

function fromMinorUnits(amount: unknown): number {
  return toNumber(amount, 0) / 100;
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

function createIdempotencyKey(prefix: string): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return `${prefix}:${crypto.randomUUID()}`;
  }
  return `${prefix}:${Date.now()}:${Math.random().toString(36).slice(2)}`;
}

function buildWriteHeaders(shopNo: string, prefix: string) {
  return {
    "x-idempotency-key": createIdempotencyKey(prefix),
    "x-shop-no": shopNo
  };
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
    attributeValues: toSaleAttrMap(raw?.spu_attrs ?? raw?.spuAttrs),
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
    salePrice: fromMinorUnits(raw.sale_price ?? raw.salePrice),
    marketPrice: fromMinorUnits(raw.market_price ?? raw.marketPrice),
    saleAttrs: toSaleAttrMap(raw.sale_attrs ?? raw.saleAttrs),
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

function buildSpuAttrs(attributeValues: Record<string, string>) {
  const attrLabelMap: Record<string, string> = {
    brand: "brand",
    material: "material",
    origin: "origin",
    craft: "craft",
    fitScene: "fit_scene",
    servicePromise: "service_promise"
  };

  return Object.entries(attributeValues)
    .map(([attrCode, value]) => ({
      attr_code: attrCode,
      attr_name: attrLabelMap[attrCode] ?? attrCode,
      scope: 1,
      value: String(value ?? "").trim()
    }))
    .filter((item) => item.value);
}

function toSaleAttrMap(items?: Array<{ attr_code?: string; attrCode?: string; value?: string }>): Record<string, string> {
  return (items ?? []).reduce<Record<string, string>>((acc, item) => {
    const key = toString(item.attr_code ?? item.attrCode).trim();
    if (!key) {
      return acc;
    }
    acc[key] = toString(item.value);
    return acc;
  }, {});
}

function buildSkuItems(payload: SellerCatalogDraftPayload) {
  return payload.skus.map((sku, index) => ({
    sku_no: sku.skuNo ?? "",
    sku_name: sku.name,
    sku_image_asset_id: Number(sku.skuImageAssetId ?? 0),
    sale_price: toMinorUnits(sku.salePrice),
    market_price: toMinorUnits(sku.marketPrice),
    sale_attrs: Object.entries(sku.saleSpecs ?? {})
      .map(([attrCode, value]) => ({
        attr_code: attrCode,
        attr_name: attrCode,
        value: String(value ?? "").trim()
      }))
      .filter((item) => item.value),
    sort_order: index,
    sku_status: 1
  }));
}

function mapAggregate(product?: RawProductAggregate) {
  const spu = normalizeSpu(product?.spu);
  return {
    spu,
    skus: (product?.skus ?? []).map((item) => normalizeSku(item))
  };
}

export async function saveSellerProductDraft(
  params: SaveSellerProductDraftInput
): Promise<{ product: { spu: SellerProductSpu; skus: Array<SellerProductSku & { saleAttrs: Record<string, string> }> } }> {
  const payload: SellerCatalogDraftPayload = {
    draft: params.draft,
    skus: params.skus
  };

  const productResponse = params.spuNo
    ? await apiClient.put<RawCreateDraftResponse>(
        "/v1/catalog/seller/products/draft",
        {
          spu_no: params.spuNo,
          expected_version: params.expectedVersion ?? 0,
          patch: {
            title: payload.draft.title,
            sub_title: payload.draft.summary,
            category_id: payload.draft.categoryId,
            brand_no: payload.draft.brandNo,
            main_image_asset_ids: payload.draft.mainImageAssetId ? [Number(payload.draft.mainImageAssetId)] : [],
            detail_image_asset_ids: (payload.draft.detailImageAssetIds ?? []).map((item) => Number(item)),
            spu_attrs: buildSpuAttrs(payload.draft.attributeValues ?? {})
          },
          update_mask: {
            paths: [
              "title",
              "sub_title",
              "category_id",
              "brand_no",
              "main_image_asset_ids",
              "detail_image_asset_ids",
              "spu_attrs"
            ]
          }
        },
        { headers: buildWriteHeaders(params.shopNo, "seller-product-draft-update") }
      )
    : await apiClient.post<RawCreateDraftResponse>(
        "/v1/catalog/seller/products/draft",
        {
          shop_no: params.shopNo,
          title: payload.draft.title,
          sub_title: payload.draft.summary,
          category_id: payload.draft.categoryId,
          brand_no: payload.draft.brandNo,
          main_image_asset_ids: payload.draft.mainImageAssetId ? [Number(payload.draft.mainImageAssetId)] : [],
          detail_image_asset_ids: (payload.draft.detailImageAssetIds ?? []).map((item) => Number(item)),
          spu_attrs: buildSpuAttrs(payload.draft.attributeValues ?? {})
        },
        { headers: buildWriteHeaders(params.shopNo, "seller-product-draft-create") }
      );

  const aggregateAfterDraft = mapAggregate(productResponse.product);
  const skuResponse = await apiClient.post<{ skus?: RawProductSku[]; spu_version?: number; spuVersion?: number }>(
    "/v1/catalog/seller/products/skus:upsert",
    {
      spu_no: aggregateAfterDraft.spu.spuNo,
      expected_version: aggregateAfterDraft.spu.version,
      items: buildSkuItems(payload),
      replace_all: true
    },
    { headers: buildWriteHeaders(params.shopNo, "seller-product-skus-upsert") }
  );

  return {
    product: {
      spu: {
        ...aggregateAfterDraft.spu,
        version: toNumber(skuResponse.spu_version ?? skuResponse.spuVersion, 0) || aggregateAfterDraft.spu.version
      },
      skus: (skuResponse.skus ?? []).map((item, index) => {
        const normalizedSku = normalizeSku(item);
        return {
          ...normalizedSku,
          saleAttrs: normalizedSku.saleAttrs ?? params.skus[index]?.saleSpecs ?? {}
        };
      })
    }
  };
}

export async function submitSellerProductReview(params: {
  shopNo: string;
  spuNo: string;
  expectedVersion: number;
  submitNote?: string;
}) {
  return apiClient.post(
    "/v1/catalog/seller/products/review:submit",
    {
      spu_no: params.spuNo,
      expected_version: params.expectedVersion,
      submit_note: params.submitNote ?? ""
    },
    {
      headers: buildWriteHeaders(params.shopNo, "seller-product-review-submit")
    }
  );
}

export function mapProductDetailToDraftPayload(detail: SellerProductDetail | ProductDetailResult): SellerCatalogDraftPayload {
  const normalized =
    "product" in detail
      ? {
          spu: {
            title: detail.product.spu.title,
            subTitle: detail.product.spu.summary,
            categoryId: detail.product.spu.categoryId,
            brandNo: detail.product.spu.brandNo,
            mainImageAssetIds: detail.product.spu.mainImageAssetIds,
            detailImageAssetIds: detail.product.spu.detailImageAssetIds,
            attributeValues: detail.product.spu.attributeValues
          },
          skus: detail.product.skus.map((sku) => ({
            skuNo: sku.skuNo,
            skuName: sku.skuName,
            skuImageAssetId: sku.skuImageAssetId,
            salePrice: sku.salePrice,
            marketPrice: sku.marketPrice,
            saleAttrs: sku.saleAttrs
          }))
        }
      : detail;
  return {
    draft: {
      title: normalized.spu?.title ?? "",
      summary: "subTitle" in (normalized.spu ?? {}) ? (normalized.spu as { subTitle?: string }).subTitle ?? "" : "",
      categoryId: normalized.spu?.categoryId ?? 0,
      brandNo: normalized.spu?.brandNo ?? "",
      mainImageAssetId: normalized.spu?.mainImageAssetIds?.[0],
      detailImageAssetIds: normalized.spu?.detailImageAssetIds ?? [],
      submitNote: "",
      attributeValues: normalized.spu?.attributeValues ?? {}
    },
    skus: normalized.skus.map((sku) => ({
      skuNo: sku.skuNo,
      name: sku.skuName,
      skuImageAssetId: sku.skuImageAssetId,
      salePrice: sku.salePrice,
      marketPrice: sku.marketPrice,
      saleSpecs: (sku as SellerProductSku & { saleAttrs?: Record<string, string> }).saleAttrs ?? {},
      initialStock: 0
    }))
  };
}

function toConsoleAggregate(detail: SellerProductDetail): ProductDetailResult["product"] {
  return {
    spu: {
      spuNo: detail.spu?.spuNo ?? "",
      shopNo: detail.spu?.shopNo ?? "",
      title: detail.spu?.title ?? "",
      summary: detail.spu?.subTitle ?? "",
      categoryId: detail.spu?.categoryId ?? 0,
      brandNo: detail.spu?.brandNo ?? "",
      status: spuStatusCodeMap[detail.spu?.spuStatus ?? "SPU_STATUS_DRAFT"] ?? "draft",
      mainImageAssetIds: detail.spu?.mainImageAssetIds ?? [],
      detailImageAssetIds: detail.spu?.detailImageAssetIds ?? [],
      attributeValues: detail.spu?.attributeValues ?? {},
      version: detail.spu?.version ?? 0
    },
    skus: detail.skus.map((sku) => ({
      skuNo: sku.skuNo,
      skuName: sku.skuName,
      skuImageAssetId: sku.skuImageAssetId,
      salePrice: sku.salePrice,
      marketPrice: sku.marketPrice,
      saleAttrs: sku.saleAttrs ?? {},
      sortOrder: sku.sortOrder
    }))
  };
}

export async function fetchMyProduct(spuNo: string): Promise<ProductDetailResult> {
  const detail = await getSellerProductDetail(spuNo);
  return {
    product: toConsoleAggregate(detail),
    review: detail.review
      ? {
          reviewStatus: spuStatusCodeMap[detail.spu?.spuStatus ?? "SPU_STATUS_DRAFT"] ?? "draft",
          rejectReasonCode: detail.review.rejectReasonCode,
          rejectComment: detail.review.rejectComment
        }
      : undefined
  };
}

export async function saveProductDraft(params: SaveSellerProductDraftInput): Promise<SaveDraftResult> {
  const result = await saveSellerProductDraft(params);
  return {
    product: {
      spu: {
        spuNo: result.product.spu.spuNo,
        shopNo: result.product.spu.shopNo,
        title: result.product.spu.title,
        summary: result.product.spu.subTitle,
        categoryId: result.product.spu.categoryId,
        brandNo: result.product.spu.brandNo,
        status: spuStatusCodeMap[result.product.spu.spuStatus] ?? "draft",
        mainImageAssetIds: result.product.spu.mainImageAssetIds,
        detailImageAssetIds: result.product.spu.detailImageAssetIds,
        attributeValues: result.product.spu.attributeValues ?? {},
        version: result.product.spu.version
      },
      skus: result.product.skus.map((sku) => ({
        skuNo: sku.skuNo,
        skuName: sku.skuName,
        skuImageAssetId: sku.skuImageAssetId,
        salePrice: sku.salePrice,
        marketPrice: sku.marketPrice,
        saleAttrs: sku.saleAttrs ?? {},
        sortOrder: sku.sortOrder
      }))
    }
  };
}

export async function submitProductReview(params: {
  shopNo: string;
  spuNo: string;
  expectedVersion: number;
  submitNote?: string;
}) {
  return submitSellerProductReview(params);
}

export async function deleteProductDraft(params: DeleteProductDraftInput) {
  return apiClient.post(
    "/v1/catalog/seller/products/draft:delete",
    {
      spu_no: params.spuNo,
      expected_version: params.expectedVersion,
      reason_code: params.reasonCode ?? "MANUAL"
    },
    {
      headers: buildWriteHeaders(params.shopNo, "seller-product-draft-delete")
    }
  );
}
