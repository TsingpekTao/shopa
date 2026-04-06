export type SpuStatusCode =
  | "SPU_STATUS_UNSPECIFIED"
  | "SPU_STATUS_DRAFT"
  | "SPU_STATUS_REVIEWING"
  | "SPU_STATUS_APPROVED"
  | "SPU_STATUS_ON_SHELF"
  | "SPU_STATUS_OFF_SHELF"
  | "SPU_STATUS_REJECTED"
  | "SPU_STATUS_FROZEN"
  | "SPU_STATUS_DELETED";

export type StockStatusCode =
  | "STOCK_STATUS_UNSPECIFIED"
  | "STOCK_STATUS_IN_STOCK"
  | "STOCK_STATUS_OUT_OF_STOCK";

export interface SellerProductSpu {
  spuNo: string;
  shopNo: string;
  title: string;
  subTitle: string;
  categoryId: number;
  brandNo: string;
  mainImageAssetIds: string[];
  detailImageAssetIds: string[];
  attributeValues?: Record<string, string>;
  spuStatus: SpuStatusCode;
  spuStockStatus: StockStatusCode;
  minSalePrice: number;
  maxSalePrice: number;
  minMarketPrice: number;
  maxMarketPrice: number;
  version: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface SellerProductSku {
  skuNo: string;
  spuNo: string;
  skuName: string;
  skuImageAssetId?: string;
  salePrice: number;
  marketPrice: number;
  saleAttrs?: Record<string, string>;
  stockStatus: StockStatusCode;
  stockVersion: number;
  sortOrder: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface SellerProductReviewInfo {
  rejectReasonCode?: string;
  rejectComment?: string;
  reviewerId?: string;
  reviewedAt?: string;
}

export interface SellerProductDetail {
  spu: SellerProductSpu | null;
  skus: SellerProductSku[];
  review?: SellerProductReviewInfo;
}

export interface ListSellerProductsParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  statuses?: SpuStatusCode[];
}

export interface ListSellerProductsResult {
  products: SellerProductSpu[];
  page: number;
  pageSize: number;
  total: number;
}

export interface DeleteProductDraftInput {
  shopNo: string;
  spuNo: string;
  expectedVersion: number;
  reasonCode?: string;
}

export interface CreateProductDraftInput {
  shopNo: string;
  title: string;
  subTitle?: string;
  categoryId?: number;
  brandNo?: string;
  mainImageAssetIds?: string[];
  detailImageAssetIds?: string[];
}

export interface CreateProductDraftResult {
  spuNo: string;
  version: number;
}

export interface SellerCatalogDraft {
  title: string;
  summary: string;
  categoryId: number;
  brandNo: string;
  mainImageAssetId?: string;
  detailImageAssetIds: string[];
  submitNote?: string;
  attributeValues: Record<string, string>;
}

export interface SellerCatalogSkuDraft {
  skuNo?: string;
  name: string;
  skuImageAssetId?: string;
  salePrice: number;
  marketPrice: number;
  saleSpecs: Record<string, string>;
  initialStock?: number;
}

export interface SellerCatalogDraftPayload {
  draft: SellerCatalogDraft;
  skus: SellerCatalogSkuDraft[];
}

export interface SaveSellerProductDraftInput {
  shopNo: string;
  draft: SellerCatalogDraft;
  skus: SellerCatalogSkuDraft[];
  spuNo?: string;
  expectedVersion?: number;
}

export type SpuStatus = "draft" | "reviewing" | "approved" | "onShelf" | "offShelf";

export const spuStatusCodeMap: Record<string, SpuStatus> = {
  SPU_STATUS_DRAFT: "draft",
  SPU_STATUS_REVIEWING: "reviewing",
  SPU_STATUS_REVIEWED: "reviewing",
  SPU_STATUS_APPROVED: "approved",
  SPU_STATUS_ON_SHELF: "onShelf",
  SPU_STATUS_OFF_SHELF: "offShelf",
  SPU_STATUS_REJECTED: "draft",
  SPU_STATUS_FROZEN: "offShelf",
  SPU_STATUS_DELETED: "offShelf"
};

export type CatalogDraft = SellerCatalogDraft;
export type CatalogSkuDraft = SellerCatalogSkuDraft;
export type CatalogDraftPayload = SellerCatalogDraftPayload;

export interface CatalogAggregateSku {
  skuNo: string;
  skuName: string;
  skuImageAssetId?: string;
  salePrice: number;
  marketPrice: number;
  saleAttrs: Record<string, string>;
  sortOrder: number;
}

export interface CatalogAggregateSpu {
  spuNo: string;
  shopNo: string;
  title: string;
  summary: string;
  categoryId: number;
  brandNo: string;
  status: SpuStatus;
  mainImageAssetIds: string[];
  detailImageAssetIds: string[];
  attributeValues: Record<string, string>;
  version: number;
}

export interface CatalogProductAggregate {
  spu: CatalogAggregateSpu;
  skus: CatalogAggregateSku[];
}

export interface CatalogReviewInfo {
  reviewStatus?: string;
  rejectReasonCode?: string;
  rejectComment?: string;
}

export interface SaveDraftResult {
  product: CatalogProductAggregate;
}

export interface ProductDetailResult {
  product: CatalogProductAggregate;
  review?: CatalogReviewInfo;
}
