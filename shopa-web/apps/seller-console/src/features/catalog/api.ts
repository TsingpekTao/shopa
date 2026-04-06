import { apiClient, ApiRequestConfig } from "@shopa/api-client";
import {
  CatalogApiProduct,
  CatalogApiResponse,
  CatalogDraftPayload,
  CatalogProduct,
  CatalogProductAggregate,
  CatalogReviewInfo,
  ProductDetailResult,
  SaveDraftResult,
  spuStatusCodeMap
} from "./types";

type CatalogApiAttribute = {
  attrCode?: string;
  attr_code?: string;
  attrName?: string;
  attr_name?: string;
  value?: string;
};

type CatalogApiSaleAttr = {
  attrCode?: string;
  attr_code?: string;
  attrName?: string;
  attr_name?: string;
  value?: string;
};

type CatalogApiSku = {
  skuNo?: string;
  sku_no?: string;
  skuName?: string;
  sku_name?: string;
  skuImageAssetId?: number | string;
  sku_image_asset_id?: number | string;
  salePrice?: number | string;
  sale_price?: number | string;
  marketPrice?: number | string;
  market_price?: number | string;
  saleAttrs?: CatalogApiSaleAttr[];
  sale_attrs?: CatalogApiSaleAttr[];
  sortOrder?: number;
  sort_order?: number;
};

type CatalogApiSpu = CatalogApiProduct & {
  shopNo?: string;
  shop_no?: string;
  subTitle?: string;
  sub_title?: string;
  categoryId?: number | string;
  category_id?: number | string;
  brandNo?: string;
  brand_no?: string;
  detailImageAssetIds?: Array<number | string>;
  detail_image_asset_ids?: Array<number | string>;
  spuAttrs?: CatalogApiAttribute[];
  spu_attrs?: CatalogApiAttribute[];
  version?: number;
};

type ProductAggregateResponse = {
  spu?: CatalogApiSpu;
  skus?: CatalogApiSku[];
};

type ReviewInfoResponse = {
  reviewStatus?: string;
  review_status?: string;
  rejectReasonCode?: string;
  reject_reason_code?: string;
  rejectComment?: string;
  reject_comment?: string;
};

type GetMyProductResponse = {
  product?: ProductAggregateResponse;
  review?: ReviewInfoResponse;
};

type CreateOrUpdateDraftResponse = {
  product?: ProductAggregateResponse;
};

type UpsertSkuDraftsResponse = {
  skus?: CatalogApiSku[];
  spuVersion?: number;
  spu_version?: number;
};

const fallbackProducts: CatalogProduct[] = [];
const DEFAULT_SHOP_NO = process.env.NEXT_PUBLIC_SELLER_SHOP_NO ?? "";

function toNumber(value?: number | string): number {
  if (typeof value === "number") {
    return value;
  }
  if (typeof value === "string") {
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : 0;
  }
  return 0;
}

function createIdempotencyKey(prefix: string): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return `${prefix}:${crypto.randomUUID()}`;
  }
  return `${prefix}:${Date.now()}:${Math.random().toString(36).slice(2)}`;
}

export function toMinorUnits(amount: number): number {
  return Math.max(0, Math.round((Number.isFinite(amount) ? amount : 0) * 100));
}

export function fromMinorUnits(amount: number): number {
  return toNumber(amount) / 100;
}

function toStringList(items?: Array<number | string>): string[] {
  return (items ?? []).map((item) => String(item)).filter(Boolean);
}

function toAttributeMap(items?: CatalogApiAttribute[]): Record<string, string> {
  return (items ?? []).reduce<Record<string, string>>((acc, item) => {
    const key = String(item.attrCode ?? item.attr_code ?? "").trim();
    if (!key) {
      return acc;
    }
    acc[key] = String(item.value ?? "");
    return acc;
  }, {});
}

function toSaleAttrMap(items?: CatalogApiSaleAttr[]): Record<string, string> {
  return (items ?? []).reduce<Record<string, string>>((acc, item) => {
    const key = String(item.attrCode ?? item.attr_code ?? "").trim();
    if (!key) {
      return acc;
    }
    acc[key] = String(item.value ?? "");
    return acc;
  }, {});
}

function buildWriteConfig(shopNo: string, prefix: string): ApiRequestConfig {
  return {
    headers: {
      "x-idempotency-key": createIdempotencyKey(prefix),
      "x-shop-no": shopNo
    }
  };
}

function mapProduct(item: CatalogApiProduct): CatalogProduct {
  const statusCode = item.spuStatusCode ?? item.spuStatus ?? item.spu_status ?? "SPU_STATUS_DRAFT";
  const mainImageAssetId = item.mainImageAssetIds?.[0] ?? item.main_image_asset_ids?.[0];
  return {
    spuNo: item.spuNo ?? item.spu_no ?? "",
    title: item.title,
    status: spuStatusCodeMap[statusCode] ?? "draft",
    salePrice: fromMinorUnits(toNumber(item.minSalePrice ?? item.min_sale_price)),
    mainImageAssetId: mainImageAssetId ? String(mainImageAssetId) : undefined
  };
}

function mapAggregate(product?: ProductAggregateResponse): CatalogProductAggregate {
  const spu = (product?.spu ?? {}) as CatalogApiSpu;
  const statusCode = spu.spuStatusCode ?? spu.spuStatus ?? spu.spu_status ?? "SPU_STATUS_DRAFT";
  return {
    spu: {
      spuNo: String(spu.spuNo ?? spu.spu_no ?? ""),
      shopNo: String(spu.shopNo ?? spu.shop_no ?? ""),
      title: String(spu.title ?? ""),
      summary: String(spu.subTitle ?? spu.sub_title ?? ""),
      categoryId: toNumber(spu.categoryId ?? spu.category_id),
      brandNo: String(spu.brandNo ?? spu.brand_no ?? ""),
      status: spuStatusCodeMap[statusCode] ?? "draft",
      mainImageAssetIds: toStringList(spu.mainImageAssetIds ?? spu.main_image_asset_ids),
      detailImageAssetIds: toStringList(spu.detailImageAssetIds ?? spu.detail_image_asset_ids),
      attributeValues: toAttributeMap(spu.spuAttrs ?? spu.spu_attrs),
      version: toNumber(spu.version)
    },
    skus: (product?.skus ?? []).map((sku) => ({
      skuNo: String(sku.skuNo ?? sku.sku_no ?? ""),
      skuName: String(sku.skuName ?? sku.sku_name ?? ""),
      skuImageAssetId: sku.skuImageAssetId ?? sku.sku_image_asset_id ? String(sku.skuImageAssetId ?? sku.sku_image_asset_id) : undefined,
      salePrice: fromMinorUnits(toNumber(sku.salePrice ?? sku.sale_price)),
      marketPrice: fromMinorUnits(toNumber(sku.marketPrice ?? sku.market_price)),
      saleAttrs: toSaleAttrMap(sku.saleAttrs ?? sku.sale_attrs),
      sortOrder: toNumber(sku.sortOrder ?? sku.sort_order)
    }))
  };
}

function mapReview(review?: ReviewInfoResponse): CatalogReviewInfo | undefined {
  if (!review) {
    return undefined;
  }
  return {
    reviewStatus: review.reviewStatus ?? review.review_status,
    rejectReasonCode: review.rejectReasonCode ?? review.reject_reason_code,
    rejectComment: review.rejectComment ?? review.reject_comment
  };
}

function buildSpuAttrs(attributeValues: Record<string, string>) {
  const attrLabelMap: Record<string, string> = {
    brand: "品牌",
    material: "材质",
    origin: "产地",
    color: "主色",
    craft: "工艺亮点",
    fitScene: "适用场景",
    servicePromise: "服务承诺"
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

function buildSkuItems(payload: CatalogDraftPayload) {
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

export async function fetchProductList(): Promise<CatalogProduct[]> {
  try {
    const response = await apiClient.get<CatalogApiResponse>("/v1/catalog/seller/products", {
      params: {
        page: 1,
        page_size: 20
      },
      silentDegraded: true
    });
    return (response.products ?? []).map(mapProduct);
  } catch (error) {
    return fallbackProducts;
  }
}

export async function fetchMyProduct(spuNo: string): Promise<ProductDetailResult> {
  const response = await apiClient.get<GetMyProductResponse>(`/v1/catalog/seller/products/${spuNo}`, {
    silentDegraded: true
  });
  return {
    product: mapAggregate(response.product),
    review: mapReview(response.review)
  };
}

export async function saveProductDraft(params: {
  shopNo: string;
  draft: CatalogDraftPayload["draft"];
  skus: CatalogDraftPayload["skus"];
  spuNo?: string;
  expectedVersion?: number;
}): Promise<SaveDraftResult> {
  const payload: CatalogDraftPayload = {
    draft: params.draft,
    skus: params.skus
  };

  let product = params.spuNo
    ? (
        await apiClient.put<CreateOrUpdateDraftResponse>(
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
          buildWriteConfig(params.shopNo, "catalog-update-draft")
        )
      ).product
    : (
        await apiClient.post<CreateOrUpdateDraftResponse>(
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
          buildWriteConfig(params.shopNo, "catalog-create-draft")
        )
      ).product;

  const aggregateAfterDraft = mapAggregate(product);
  const skuResponse = await apiClient.post<UpsertSkuDraftsResponse>(
    "/v1/catalog/seller/products/skus:upsert",
    {
      spu_no: aggregateAfterDraft.spu.spuNo,
      expected_version: aggregateAfterDraft.spu.version,
      items: buildSkuItems(payload),
      replace_all: true
    },
    buildWriteConfig(params.shopNo, "catalog-upsert-skus")
  );

  return {
    product: {
      spu: {
        ...aggregateAfterDraft.spu,
        version: toNumber(skuResponse.spuVersion ?? skuResponse.spu_version) || aggregateAfterDraft.spu.version
      },
      skus: (skuResponse.skus ?? []).map((sku, index) => ({
        skuNo: String(sku.skuNo ?? sku.sku_no ?? ""),
        skuName: String(sku.skuName ?? sku.sku_name ?? payload.skus[index]?.name ?? ""),
        skuImageAssetId: sku.skuImageAssetId ?? sku.sku_image_asset_id ? String(sku.skuImageAssetId ?? sku.sku_image_asset_id) : undefined,
        salePrice: fromMinorUnits(toNumber(sku.salePrice ?? sku.sale_price)),
        marketPrice: fromMinorUnits(toNumber(sku.marketPrice ?? sku.market_price)),
        saleAttrs: toSaleAttrMap(sku.saleAttrs ?? sku.sale_attrs),
        sortOrder: toNumber(sku.sortOrder ?? sku.sort_order)
      }))
    }
  };
}

export async function submitDraft(payload: CatalogDraftPayload): Promise<{ success: boolean }> {
  await saveProductDraft({
    shopNo: DEFAULT_SHOP_NO,
    draft: payload.draft,
    skus: payload.skus
  });
  return { success: true };
}

export async function submitProductReview(params: {
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
    buildWriteConfig(params.shopNo, "catalog-submit-review")
  );
}

export function mapProductDetailToDraftPayload(detail: ProductDetailResult): CatalogDraftPayload {
  return {
    draft: {
      title: detail.product.spu.title,
      summary: detail.product.spu.summary,
      categoryId: detail.product.spu.categoryId,
      brandNo: detail.product.spu.brandNo,
      mainImageAssetId: detail.product.spu.mainImageAssetIds[0],
      detailImageAssetIds: detail.product.spu.detailImageAssetIds,
      submitNote: "",
      attributeValues: detail.product.spu.attributeValues
    },
    skus: detail.product.skus.map((sku) => ({
      skuNo: sku.skuNo,
      name: sku.skuName,
      skuImageAssetId: sku.skuImageAssetId,
      salePrice: sku.salePrice,
      marketPrice: sku.marketPrice,
      saleSpecs: sku.saleAttrs,
      initialStock: 0
    }))
  };
}


