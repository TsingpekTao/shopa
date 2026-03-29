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
